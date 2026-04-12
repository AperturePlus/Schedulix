# 阶段六：Serverless 架构

## 学习目标

将调度系统封装为 HTTP API，学习 Serverless 架构的核心概念。

## 前置知识

- 阶段五完成
- 了解 HTTP 协议基础（GET/POST、状态码、JSON）

## 核心概念

### 1. Serverless 核心思想

传统服务器：你管理服务器，7×24 运行，即使没有请求也占用资源。

Serverless：
- **无服务器管理**：不需要关心服务器的运维
- **按需执行**：有请求时启动函数，无请求时释放资源
- **自动扩缩容**：根据负载自动调整实例数量
- **缩容到零**：完全没有请求时，实例数降为 0

### 2. 冷启动 vs 热启动

```
冷启动（首次请求）：
请求到达 → 创建实例 → 初始化环境 → 执行函数 → 返回响应
           ├── 冷启动延迟 ──┤

热启动（后续请求）：
请求到达 → 执行函数 → 返回响应
           （实例已存在，无需初始化）
```

在 Schedulix 中用 `time.Sleep` 模拟冷启动延迟。

### 3. net/http 标准库

Go 的 HTTP 服务非常简洁：

```go
// 注册路由
mux := http.NewServeMux()
mux.HandleFunc("POST /api/v1/tasks", submitTaskHandler)
mux.HandleFunc("GET /api/v1/tasks/{id}", getTaskHandler)

// 启动服务
http.ListenAndServe(":8080", mux)
```

**Handler 函数签名**：
```go
func handler(w http.ResponseWriter, r *http.Request) {
    // w: 写响应
    // r: 读请求
}
```

### 4. 请求参数验证

```go
func submitTask(w http.ResponseWriter, r *http.Request) {
    var task model.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
        return
    }
    
    // 验证必填字段
    if task.ID == "" {
        http.Error(w, `{"error": "task ID is required"}`, http.StatusBadRequest)
        return
    }
    
    // 验证资源需求
    if task.Resource.Memory <= 0 {
        http.Error(w, `{"error": "memory must be positive"}`, http.StatusBadRequest)
        return
    }
    
    // 处理请求...
    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(map[string]string{"id": task.ID, "status": "accepted"})
}
```

### 5. 依赖注入

Handler 不直接创建依赖，而是通过构造函数注入：

```go
// 好的做法：依赖注入
handler := NewHandler(cluster, scheduler, queue)

// 不好的做法：在 handler 内部创建
func handler(w http.ResponseWriter, r *http.Request) {
    cluster := model.NewCluster(100) // ❌ 每次请求都创建新集群
}
```

### 6. 自动扩缩容模拟

```go
type Scaler struct {
    activeInstances atomic.Int64
    coldStartDelay  time.Duration
}

func (s *Scaler) OnRequest() bool {
    isCold := s.activeInstances.Load() == 0
    if isCold {
        time.Sleep(s.coldStartDelay) // 模拟冷启动
    }
    s.activeInstances.Add(1)
    return isCold
}

func (s *Scaler) OnRequestDone() {
    s.activeInstances.Add(-1)
    // 如果降到 0，启动延迟缩容定时器
}
```

## API 端点一览

| 方法 | 路径 | 功能 |
|------|------|------|
| POST | /api/v1/tasks | 提交任务 |
| GET | /api/v1/tasks/{id} | 查询任务状态 |
| GET | /api/v1/cluster/status | 查询集群状态 |
| GET | /api/v1/cluster/nodes | 查询节点列表 |
| POST | /api/v1/simulator/start | 启动事件模拟 |
| GET | /api/v1/metrics | 获取监控指标 |
| GET | /api/v1/metrics/export | 导出历史指标 |

## 练习任务

1. 打开 `pkg/gateway/handler.go`，实现所有 HTTP handler
2. 打开 `pkg/gateway/router.go`，注册路由
3. 打开 `pkg/gateway/scaler.go`，实现冷/热启动和扩缩容
4. 打开 `cmd/server/main.go`，组装所有组件并启动服务

## 验证

```bash
go test ./pkg/gateway/...

# 手动测试（启动服务后）
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"id":"task-1","resource":{"compute_power":10,"memory":1024},"priority":5}'
```
