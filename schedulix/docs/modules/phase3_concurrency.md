# 阶段三（扩展）：Go 并发模式大全

## 学习目标

掌握 Go 并发编程的 7 种核心模式，从基础到高级。

## 前置知识

- 阶段一完成
- 了解线程和锁的基本概念

## 七种并发模式

### 模式 1：Fan-Out / Fan-In（扇出/扇入）

```
一个生产者，多个消费者并行处理，结果汇聚。

场景：检查 10000 个节点的健康状态
串行：10000 × 10ms = 100s
Fan-Out 100 workers：10000 / 100 × 10ms = 1s

         ┌→ worker 1 →┐
input ──→├→ worker 2 →├──→ output
         └→ worker 3 →┘
```

### 模式 2：Worker Pool（工作池）

```
固定数量的长期运行 worker，复用处理任务。

任务 → [buffered channel] → worker 1 (循环处理)
                          → worker 2 (循环处理)
                          → worker 3 (循环处理)

vs 每个任务一个 goroutine：
任务 1 → goroutine 1 (用完即弃)
任务 2 → goroutine 2 (用完即弃)
...
任务 100000 → goroutine 100000 (内存爆炸)
```

### 模式 3：Pipeline（并发管道）

```
多个阶段串联，每个阶段是一个 goroutine。

[生成] ──ch1──→ [过滤] ──ch2──→ [转换] ──ch3──→ [输出]
  g1              g2              g3              g4

每个阶段独立运行，通过 channel 传递数据。
关闭传播：ch1 关闭 → g2 退出 → ch2 关闭 → g3 退出 → ...
```

### 模式 4：Select 多路复用

```go
select {
case msg := <-ch1:
    // ch1 先就绪
case msg := <-ch2:
    // ch2 先就绪
case <-time.After(5 * time.Second):
    // 超时
case <-ctx.Done():
    // 取消
}
```

### 模式 5：sync.Once / sync.Pool

```go
// Once：只执行一次（线程安全的单例）
var once sync.Once
var db *Database
once.Do(func() { db = connectDB() })

// Pool：对象复用，减少 GC
pool := sync.Pool{New: func() any { return &Buffer{} }}
buf := pool.Get().(*Buffer)
defer pool.Put(buf)
```

### 模式 6：Semaphore（信号量）

```
限制并发数。用 buffered channel 实现。

sem := make(chan struct{}, 3) // 最多 3 个并发

sem <- struct{}{}  // 获取（满了就阻塞）
defer func() { <-sem }()  // 释放
```

### 模式 7：ErrGroup（错误组）

```
并发执行多个操作，任一失败则取消其余。

g := NewErrGroup(ctx)
g.Go(func(ctx) error { return checkNode1(ctx) })
g.Go(func(ctx) error { return checkNode2(ctx) })
g.Go(func(ctx) error { return checkNode3(ctx) })
err := g.Wait() // 返回第一个错误，其余被取消
```

## 竞争检测

```bash
# 必须用 -race 运行所有并发测试
go test -race ./pkg/concurrency/...
go test -race ./pkg/scheduler/...
```

`-race` 会检测数据竞争（两个 goroutine 同时读写同一变量且至少一个是写）。
生产代码中的数据竞争 = bug，必须修复。

## 练习任务

1. `pkg/concurrency/patterns.go` — 实现所有 7 种模式
2. 在 `pkg/scheduler/concurrent.go` 中使用 WorkerPool 替代手动 goroutine 管理
3. 在 `pkg/simulator/engine.go` 中使用 FanOut 并行检查节点
4. 在 `pkg/orchestrator/replicaset.go` 中使用 ErrGroup 并发创建 Pod

## 验证

```bash
go test -race -v ./pkg/concurrency/...
```
