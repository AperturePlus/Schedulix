# 阶段三：并发编程

## 学习目标

利用 Go 的 goroutine、channel 和 sync 包实现并发调度，深入理解 Go 的并发模型。

## 前置知识

- 阶段二完成
- 了解并发与并行的区别

## 核心概念

### 1. Goroutine

Go 的轻量级线程，由 Go 运行时调度（不是操作系统线程）。

```go
// 启动一个 goroutine
go func() {
    fmt.Println("我在另一个 goroutine 中运行")
}()
```

**关键**：goroutine 非常轻量（初始栈只有几 KB），可以轻松创建数千个。

### 2. Channel（通道）

goroutine 之间的通信管道。

```go
// 无缓冲 channel：发送和接收必须同时就绪
ch := make(chan *Task)

// 有缓冲 channel：缓冲区满之前发送不阻塞
ch := make(chan *Task, 1000)
```

**生产者-消费者模式**：

```
生产者 goroutine ──→ [channel 缓冲区] ──→ 消费者 goroutine
生产者 goroutine ──→ [  最多 1000 个  ] ──→ 消费者 goroutine
生产者 goroutine ──→ [    任务       ] ──→ 消费者 goroutine
```

```go
// 生产者
func producer(ch chan<- *Task, tasks []*Task) {
    for _, t := range tasks {
        ch <- t  // 发送到 channel
    }
    close(ch)  // 发送完毕，关闭 channel
}

// 消费者
func consumer(ch <-chan *Task) {
    for task := range ch {  // channel 关闭后自动退出循环
        schedule(task)
    }
}
```

### 3. sync.Mutex / sync.RWMutex

保护共享数据，防止数据竞争。

```go
var mu sync.Mutex

// 互斥锁：同一时间只有一个 goroutine 能进入
mu.Lock()
// 修改共享数据...
mu.Unlock()

// 读写锁：允许多个读，但写时独占
var rwmu sync.RWMutex
rwmu.RLock()   // 读锁，多个 goroutine 可同时持有
rwmu.RUnlock()
rwmu.Lock()    // 写锁，独占
rwmu.Unlock()
```

**何时用 RWMutex？** 读操作远多于写操作时。Schedulix 中 Cluster 的查询（读）远多于调度分配（写）。

### 4. sync.WaitGroup

等待一组 goroutine 全部完成。

```go
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        // 做一些工作...
    }()
}

wg.Wait()  // 阻塞直到所有 goroutine 调用了 Done()
```

### 5. context.Context

传播取消信号和超时控制。

```go
// 创建带超时的 context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case ch <- task:
    // 发送成功
case <-ctx.Done():
    // 超时或被取消
    return ctx.Err()
}
```

### 6. 数据竞争检测

Go 内置了竞争检测器：

```bash
go test -race ./pkg/scheduler/...
```

如果存在数据竞争，会输出详细的竞争位置和 goroutine 栈信息。**阶段三的所有测试都应该加 `-race` 标志。**

## 并发调度架构

```
HTTP Handler ──┐
HTTP Handler ──┤──→ [buffered channel (1000)] ──→ Worker goroutine 1
HTTP Handler ──┘                                ──→ Worker goroutine 2
                                                ──→ Worker goroutine 3
                                                         │
                                                    ┌────┴────┐
                                                    │ Cluster  │
                                                    │ (RWMutex)│
                                                    └──────────┘
```

**关键并发安全点**：
- Worker 调度时需要 Lock Cluster（写锁），因为要修改节点资源
- 多个 Worker 不能同时分配同一个节点的资源（会超额分配）

## 练习任务

1. 打开 `pkg/scheduler/concurrent.go`
2. 实现 `NewConcurrentScheduler`：创建 buffered channel
3. 实现 `Submit`：使用 select + ctx.Done() 发送任务
4. 实现 `StartWorkers`：启动 n 个 worker goroutine，使用 WaitGroup
5. 确保所有测试通过 `go test -race`

## 验证

```bash
go test -race ./pkg/scheduler/...
```
