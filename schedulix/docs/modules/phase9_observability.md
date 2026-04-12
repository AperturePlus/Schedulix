# 阶段九：可观测性（Observability）

## 学习目标

实现结构化日志、分布式追踪和健康检查，掌握生产系统可观测性的三大支柱。

## 前置知识

- 阶段一完成（数据模型）
- 了解 JSON 格式

## 核心概念

### 可观测性三大支柱

```
可观测性（Observability）
├── 日志（Logging）    — 发生了什么？
├── 指标（Metrics）    — 系统状态如何？（阶段八已实现）
└── 追踪（Tracing）    — 请求经过了哪些组件？
```

加上健康检查（Health Check）— 系统还活着吗？

### 1. 结构化日志 vs fmt.Println

```
// ❌ 非结构化日志 — 人能读，机器难解析
fmt.Println("2024-01-01 task-123 scheduled to node-456")

// ✅ 结构化日志 — 人和机器都能读
{"timestamp":"2024-01-01T00:00:00Z","level":"INFO","component":"scheduler","message":"task scheduled","fields":{"task_id":"task-123","node_id":"node-456","latency_ms":12}}
```

结构化日志可以用 `jq`、`grep` 等工具高效查询：
```bash
# 查找所有调度失败的日志
cat app.log | jq 'select(.level == "ERROR" and .component == "scheduler")'
```

### 2. 分布式追踪

一个请求经过多个组件，追踪记录每段的耗时：

```
请求: POST /api/v1/tasks
│
├── [Span] gateway.SubmitTask (2ms)
│   ├── [Span] queue.Enqueue (0.1ms)
│   └── [Span] scheduler.Schedule (1.5ms)
│       ├── [Span] cluster.GetAvailableNodes (0.3ms)
│       └── [Span] node.AllocateTask (0.2ms)
│
总耗时: 2ms
瓶颈: scheduler.Schedule (75%)
```

### 3. 健康检查

```go
// 注册检查项
checker.Register("cluster", func() (HealthStatus, string) {
    idle := cluster.GetAvailableNodes(NodeStatusIdle)
    if len(idle) == 0 {
        return HealthStatusUnhealthy, "no idle nodes"
    }
    if len(idle) < 10 {
        return HealthStatusDegraded, fmt.Sprintf("only %d idle nodes", len(idle))
    }
    return HealthStatusHealthy, "ok"
})

// GET /health 返回
{
  "status": "degraded",
  "components": [
    {"name": "cluster", "status": "degraded", "message": "only 5 idle nodes"},
    {"name": "queue", "status": "healthy", "message": "ok"}
  ]
}
```

## 练习任务

1. 打开 `pkg/observability/logger.go`，实现结构化日志器
2. 打开 `pkg/observability/tracer.go`，实现简化版追踪器
3. 打开 `pkg/observability/healthcheck.go`，实现健康检查器
4. 在 Gateway 中集成：添加 `GET /health` 端点
5. 在 Scheduler、Simulator、Recovery 中添加日志和追踪调用

## 验证

```bash
go test ./pkg/observability/...
```
