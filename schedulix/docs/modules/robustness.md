# 鲁棒性编程指南

## 核心原则

**一切皆不可靠。** 合格的程序员假设：
- 所有输入都可能是垃圾
- 所有指针都可能是 nil
- 所有网络调用都可能失败
- 所有并发操作都可能竞争
- 所有不变量都可能被破坏

## Schedulix 中的鲁棒性模式

### 1. 不信任输入（Never Trust Input）

每个公开方法的第一件事：验证参数。

```go
// ❌ 脆弱代码
func (n *GPU_Node) CanAccept(req ResourceRequirement) bool {
    return n.ComputePower >= req.ComputePower  // n 可能是 nil！
}

// ✅ 鲁棒代码
func (n *GPU_Node) CanAccept(req ResourceRequirement) bool {
    if n == nil {
        return false
    }
    if n.Status == NodeStatusOffline {
        return false
    }
    if req.ComputePower < 0 || req.Memory < 0 {
        return false  // 不信任调用方
    }
    // ...
}
```

### 2. 明确的错误类型（Sentinel Errors）

为每种失败模式定义独立的错误，调用方可以精确判断：

```go
var (
    ErrNoAvailableNode = errors.New("no available node")
    ErrNilTask         = errors.New("task is nil")
    ErrQueueFull       = errors.New("queue is full")
)

// 调用方
if errors.Is(err, ErrNoAvailableNode) {
    // 可恢复：等待重试
} else if errors.Is(err, ErrNilTask) {
    // 不可恢复：调用方 bug
}
```

### 3. 错误包装（Error Wrapping）

用 `%w` 包装错误，保留上下文但允许 `errors.Is` 匹配：

```go
return fmt.Errorf("schedule task %s: %w", task.ID, ErrNoAvailableNode)
// 调用方仍然可以 errors.Is(err, ErrNoAvailableNode)
// 但错误消息包含了任务 ID 上下文
```

### 4. 优雅降级（Graceful Degradation）

检测到异常时，返回安全的默认值而非 panic：

```go
// AvailableMemory：即使内部状态异常也不返回负值
func (n *GPU_Node) AvailableMemory() int {
    avail := n.MemoryTotal - n.MemoryUsed
    if avail < 0 {
        return 0  // 安全降级
    }
    return avail
}

// EventConfig：输入为空时返回默认配置
func ParseConfig(data []byte) (*EventConfig, error) {
    if len(data) == 0 {
        return DefaultEventConfig(), fmt.Errorf("empty config, using defaults")
    }
    // ...
}
```

### 5. 幂等操作（Idempotency）

同一操作执行多次，效果与执行一次相同：

```go
// 重复分配同一任务不应重复扣资源
func (n *GPU_Node) AllocateTask(taskID string, req ResourceRequirement) error {
    for _, id := range n.AssignedTasks {
        if id == taskID {
            return nil  // 已分配，幂等返回
        }
    }
    // 首次分配...
}

// 重复关闭 channel 不 panic
func (cs *ConcurrentScheduler) Stop() {
    cs.closeMu.Lock()
    defer cs.closeMu.Unlock()
    if !cs.closed {
        cs.closed = true
        close(cs.taskChan)
    }
}
```

### 6. 任务不丢失（No Task Loss）

调度失败时，任务必须回到队列或返回给调用方：

```go
func (s *Scheduler) ScheduleNext() (*model.Task, error) {
    task, err := s.queue.Dequeue()
    if err != nil {
        return nil, err
    }
    
    nodeID, err := s.strategy.Schedule(task, s.cluster)
    if err != nil {
        // 调度失败 → 任务回到队列
        if requeueErr := s.queue.Enqueue(task); requeueErr != nil {
            // 入队也失败 → 返回任务给调用方（绝不丢失）
            return task, fmt.Errorf("schedule failed and requeue failed: %w", requeueErr)
        }
        return nil, err
    }
    // ...
}
```

### 7. Panic 隔离（Panic Recovery）

单个任务/handler 的 panic 不能杀死整个系统：

```go
// Worker goroutine 内部
func safeSchedule(task *model.Task) (result ScheduleResult) {
    defer func() {
        if r := recover(); r != nil {
            result.Err = fmt.Errorf("panic: %v", r)
        }
    }()
    // 执行调度...
    return
}

// 观察者通知
for _, obs := range observers {
    func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("observer panic: %v", r)
            }
        }()
        obs.OnStateChange(id, oldState, newState)
    }()
}
```

### 8. 防御性并发（Defensive Concurrency）

```go
// 总是用 defer 释放锁
c.mu.Lock()
defer c.mu.Unlock()

// 返回副本而非内部引用
func (es *EventSimulator) GetEventLog() []*FaultEvent {
    es.logMu.RLock()
    defer es.logMu.RUnlock()
    result := make([]*FaultEvent, len(es.eventLog))
    copy(result, es.eventLog)
    return result
}

// 检查 channel 是否已关闭再发送
cs.closeMu.Lock()
if cs.closed {
    cs.closeMu.Unlock()
    return ErrChannelClosed
}
cs.closeMu.Unlock()
```

### 9. 容量限制（Capacity Limits）

防止资源耗尽：

```go
// 队列最大容量
const DefaultMaxQueueSize = 100000

// 环形缓冲区固定大小
const DefaultBufferSize = 1000

// HTTP body 大小限制
r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB
```

### 10. 索引一致性修复（Self-Healing）

当检测到内部索引不一致时，自动修复而非 panic：

```go
func (c *Cluster) GetAvailableNodes(status NodeStatus) []*GPU_Node {
    ids := c.statusIndex[status]
    result := make([]*GPU_Node, 0, len(ids))
    for _, id := range ids {
        node, ok := c.Nodes[id]
        if !ok {
            // 索引损坏！记录警告，跳过，不 panic
            log.Printf("WARNING: node %s in statusIndex but not in Nodes map", id)
            continue
        }
        result = append(result, node)
    }
    return result
}
```

## 检查清单

每个 TODO(learner) 实现前，问自己：

- [ ] 所有指针参数检查了 nil 吗？
- [ ] 所有字符串参数检查了空值吗？
- [ ] 所有数值参数检查了范围吗？
- [ ] 除法运算检查了除数为 0 吗？
- [ ] 切片操作检查了空切片吗？
- [ ] map 查找检查了 key 不存在吗？
- [ ] 错误返回了明确的类型吗？
- [ ] 错误包含了足够的上下文吗？
- [ ] 并发操作加锁了吗？用了 defer 释放吗？
- [ ] 返回的是副本还是内部引用？
- [ ] 操作是幂等的吗？
- [ ] 失败时资源回滚了吗？
- [ ] panic 被 recover 了吗？
