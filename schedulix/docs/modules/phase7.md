# 阶段七：容器操作

## 学习目标

模拟容器的生命周期管理，掌握状态机设计和资源隔离概念。

## 前置知识

- 阶段一完成（数据模型）
- 了解容器的基本概念（Docker）

## 核心概念

### 1. 有限状态机（Finite State Machine, FSM）

容器的生命周期是一个严格的状态机：

```
Created ──→ Running ──→ Stopped ──→ Destroyed
   │            │           │
   ✗            ✗           ✗
   │            │           │
   不能直接     不能直接     不能回到
   Destroyed    Created     Running
```

**合法转换**：
- Created → Running（启动）
- Running → Stopped（停止）
- Stopped → Destroyed（销毁）

**非法转换**（返回错误）：
- Created → Stopped ❌
- Created → Destroyed ❌
- Running → Created ❌
- Stopped → Running ❌

实现方式：用 map 定义合法转换表

```go
var validTransitions = map[ContainerState][]ContainerState{
    ContainerCreated: {ContainerRunning},
    ContainerRunning: {ContainerStopped},
    ContainerStopped: {ContainerDestroyed},
}
```

### 2. 资源配额

每个容器声明自己需要的资源（CPU 份额、内存限制）。宿主节点上所有容器的资源总和不能超过节点总资源。

```
GPU_Node (总内存 80GB)
├── Container A: 20GB
├── Container B: 30GB
├── Container C: 25GB
└── 剩余: 5GB

新容器 D 需要 10GB → 拒绝！(5GB < 10GB)
新容器 E 需要 3GB  → 允许！(5GB >= 3GB)
```

### 3. 观察者模式（容器版）

容器状态变更时通知订阅者：

```go
type ContainerLifecycle interface {
    OnStateChange(containerID string, oldState, newState ContainerState)
}

// 使用
runtime.RegisterObserver(metricsCollector)
runtime.RegisterObserver(logger)

// 状态变更时
for _, obs := range observers {
    obs.OnStateChange(container.ID, oldState, newState)
}
```

## 练习任务

1. 打开 `pkg/container/lifecycle.go`，实现 `IsValidTransition` 和 `TransitionState`
2. 打开 `pkg/container/runtime.go`，实现 `CreateContainer`、`StartContainer`、`StopContainer`、`DestroyContainer`

## 验证

```bash
go test ./pkg/container/...
```
