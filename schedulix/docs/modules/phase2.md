# 阶段二：任务队列与基础调度

## 学习目标

实现优先级队列和三种经典调度算法，掌握 Go 接口设计和策略模式。

## 前置知识

- 阶段一完成
- 了解堆（Heap）数据结构的基本概念

## 核心概念

### 1. container/heap 接口

Go 标准库提供了堆操作，但需要你实现底层数据结构。需要实现 5 个方法：

```go
type Interface interface {
    sort.Interface        // Len, Less, Swap
    Push(x any)
    Pop() any
}
```

**关键理解**：`heap.Push` 和 `heap.Pop` 是库函数，它们会调用你实现的 `Push`/`Pop` 方法，并自动维护堆的性质。你的 `Push` 只需要 append，`Pop` 只需要取最后一个元素。

```go
// 你的 Push：只管 append
func (h *taskHeap) Push(x any) {
    *h = append(*h, x.(*model.Task))
}

// 你的 Pop：只管取最后一个
func (h *taskHeap) Pop() any {
    old := *h
    n := len(old)
    item := old[n-1]
    *h = old[:n-1]
    return item
}

// 使用时调用库函数（不是你的方法）：
heap.Push(&h, task)    // 库函数会调用你的 Push + 上浮调整
task := heap.Pop(&h)   // 库函数会调用你的 Pop + 下沉调整
```

### 2. 策略模式（Strategy Pattern）

将算法封装为独立对象，通过接口实现可互换：

```
┌─────────────┐     ┌──────────────────┐
│  Scheduler   │────→│ ScheduleStrategy │ (接口)
└─────────────┘     └──────────────────┘
                           ▲
                    ┌──────┼──────┐
                    │      │      │
              FirstFit  BestFit  RoundRobin
```

Go 中用 interface 实现：

```go
type ScheduleStrategy interface {
    Schedule(task *Task, cluster *Cluster) (nodeID string, err error)
    Name() string
}
```

### 3. 调度算法对比

#### First-Fit（首次适应）

```
节点:  [A:空闲80MB] [B:空闲200MB] [C:空闲50MB]
任务:  需要 60MB

→ 从头遍历，A 满足（80 >= 60），选 A
→ 不看 B 和 C
```

**优点**：快，O(n) 最坏但通常很快找到
**缺点**：前面的节点容易过载

#### Best-Fit（最佳适应）

```
节点:  [A:空闲80MB] [B:空闲200MB] [C:空闲65MB]
任务:  需要 60MB

→ 遍历所有节点
→ A: 80-60=20 剩余
→ B: 200-60=140 剩余
→ C: 65-60=5 剩余 ← 最小剩余
→ 选 C
```

**优点**：资源利用率高，减少碎片
**缺点**：必须遍历所有节点

#### Round-Robin（轮询）

```
节点:  [A] [B] [C] [D]
游标:  ^

任务1 → A, 游标移到 B
任务2 → B, 游标移到 C
任务3 → C, 游标移到 D
任务4 → D, 游标移到 A（循环）
```

**优点**：负载均匀，公平
**缺点**：不考虑节点差异

## 练习任务

1. 打开 `pkg/queue/priority_queue.go`，实现 `taskHeap` 的 5 个方法和 `TaskQueue` 的所有方法
2. 打开 `pkg/scheduler/firstfit.go`，实现 First-Fit 算法
3. 打开 `pkg/scheduler/bestfit.go`，实现 Best-Fit 算法
4. 打开 `pkg/scheduler/roundrobin.go`，实现 Round-Robin 算法
5. 打开 `pkg/scheduler/scheduler.go`，实现 `ScheduleNext` 和 `ScheduleTask`

## 验证

```bash
go test ./pkg/queue/...
go test ./pkg/scheduler/...
```
