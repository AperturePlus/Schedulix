# 阶段零（可选）：Go 函数式编程

## 学习目标

掌握 Go 中函数作为一等公民的用法：函数类型、高阶函数、闭包、函数组合、函数选项模式和中间件模式。

## 前置知识

- Go 基础语法（变量、if/for、结构体）

## 为什么需要学这些？

Schedulix 项目中到处都在用这些模式：
- **调度策略接口** → 函数类型的面向对象替代
- **事件处理器** → 回调函数
- **并发 worker** → goroutine + 闭包
- **HTTP 中间件** → 函数包裹函数
- **配置选项** → 函数选项模式

## 核心概念

### 1. 函数类型

给函数签名起一个名字，像类型一样使用：

```go
// 定义函数类型
type NodePredicate func(node *GPU_Node) bool

// 使用函数类型作为参数
func FilterNodes(nodes []*GPU_Node, predicate NodePredicate) []*GPU_Node {
    result := make([]*GPU_Node, 0)
    for _, n := range nodes {
        if predicate(n) {
            result = append(result, n)
        }
    }
    return result
}

// 调用时传入匿名函数
idle := FilterNodes(nodes, func(n *GPU_Node) bool {
    return n.Status == NodeStatusIdle
})
```

### 2. 高阶函数

接收函数作为参数，或返回函数的函数：

```go
// 接收函数
func MapNodeScores(nodes []*GPU_Node, scorer func(*GPU_Node) float64) map[string]float64

// 返回函数
func NegatePredicate(p NodePredicate) NodePredicate {
    return func(n *GPU_Node) bool {
        return !p(n)
    }
}
```

### 3. 闭包（Closure）

函数 + 它捕获的外部变量：

```go
func MakeCounter() (increment func(), get func() int) {
    count := 0  // 被两个闭包共享
    increment = func() { count++ }
    get = func() int { return count }
    return
}

inc, get := MakeCounter()
inc(); inc(); inc()
fmt.Println(get()) // 3
```

闭包的常见陷阱：

```go
// ❌ 错误：循环变量被共享
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i) // 可能全部打印 5
    }()
}

// ✅ 正确：通过参数传递
for i := 0; i < 5; i++ {
    go func(n int) {
        fmt.Println(n) // 打印 0,1,2,3,4
    }(i)
}
```

### 4. 函数选项模式

```go
// 不用函数选项：参数太多，难以维护
NewScheduler(strategy, queue, cluster, 3, 100*time.Millisecond, 5*time.Second, 4, true, false)

// 用函数选项：清晰、可扩展
NewScheduler(strategy, queue, cluster,
    WithMaxRetries(3),
    WithTimeout(5 * time.Second),
    WithWorkerCount(4),
    WithMetrics(true),
)
```

### 5. 中间件模式

```go
// 中间件 = 包裹 handler 的函数
type Middleware func(http.Handler) http.Handler

// 日志中间件
func LoggingMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            next.ServeHTTP(w, r)  // 调用下一个
            log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
        })
    }
}

// 串联
handler := Chain(
    RecoveryMiddleware(),   // 最外层：捕获 panic
    LoggingMiddleware(),    // 记录日志
    RateLimitMiddleware(100, time.Second), // 限流
)(myHandler)
```

### 6. Pipeline 管道

```go
// 链式调用，数据依次流过每个步骤
result := NewNodePipeline().
    Filter(isIdle).
    Filter(hasEnoughMemory).
    SortBy(byComputePowerDesc).
    Limit(10).
    Execute(allNodes)
```

## 练习任务

1. `pkg/functools/pipeline.go` — 实现 Filter、Map、Reduce、Compose、Pipeline
2. `pkg/functools/options.go` — 实现函数选项模式、闭包（Counter、RateLimiter、Retrier）
3. `pkg/functools/middleware.go` — 实现 HTTP 中间件和调度中间件

## 验证

```bash
go test ./pkg/functools/...
```
