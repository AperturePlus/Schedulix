# 阶段十一：容器编排（Container Orchestration）

## 学习目标

模拟 Kubernetes 核心编排概念，掌握 Pod 调度、副本管理、滚动更新、服务发现和自愈机制。

## 前置知识

- 阶段七完成（容器生命周期）
- 阶段三完成（并发编程）
- 了解 K8s 的基本概念（Pod、Deployment、Service）

## 核心概念

### K8s 架构概览

```
用户 → kubectl → API Server → etcd（存储）
                     ↓
              Controller Manager（控制循环）
                     ↓
                Scheduler（调度）
                     ↓
              Kubelet（节点代理）→ 容器运行时
```

Schedulix 模拟的部分：
- **Scheduler** → `PodScheduler`（过滤 → 打分 → 绑定）
- **Controller Manager** → `ReplicaSetController` + `DeploymentController`
- **Service/kube-proxy** → `ServiceController`

### 1. Pod — 最小调度单元

```
Pod "web-abc123"
├── 状态: Running
├── 节点: node-0042
├── 标签: {app: web, version: v2}
├── 资源: {cpu: 100m, memory: 256Mi}
├── 探针:
│   ├── Liveness: 每 10s 检查，失败 3 次 → 重启
│   └── Readiness: 每 5s 检查，失败 2 次 → 从 Service 移除
└── 重启策略: Always
```

### 2. 控制循环（Reconciliation Loop）

K8s 的灵魂。所有控制器都遵循同一模式：

```
┌─────────────────────────────────────┐
│         Reconciliation Loop          │
│                                     │
│  1. 观察当前状态（Observe）          │
│  2. 比较期望状态（Compare）          │
│  3. 执行动作收敛（Act）              │
│  4. 重复                            │
│                                     │
│  期望: 3 个 Pod                     │
│  当前: 2 个 Pod                     │
│  动作: 创建 1 个 Pod                │
└─────────────────────────────────────┘
```

```go
func (rc *ReplicaSetController) Reconcile(rsID string) error {
    rs := rc.rsets[rsID]
    currentPods := rc.getMatchingPods(rs.Selector)
    
    running := countByPhase(currentPods, PodRunning, PodPending)
    desired := rs.Replicas
    
    if running < desired {
        // 扩容
        for i := 0; i < desired - running; i++ {
            pod := createPodFromTemplate(rs.Template)
            rc.scheduler.SchedulePod(pod)
        }
    } else if running > desired {
        // 缩容
        toDelete := selectPodsToDelete(currentPods, running - desired)
        for _, pod := range toDelete {
            pod.Phase = PodSucceeded
        }
    }
}
```

### 3. 滚动更新

```
v1: [●] [●] [●]     MaxUnavailable=1, MaxSurge=1
v2: (空)

步骤 1: 创建 v2 Pod（surge）
v1: [●] [●] [●]
v2: [○]              总数=4 (3+1 surge)

步骤 2: v2 就绪，终止 v1
v1: [●] [●] [✗]
v2: [●]              可用=3 (2+1 ≥ 3-1)

步骤 3-6: 重复...
v1: (空)
v2: [●] [●] [●]     完成！
```

约束条件：
- 任何时刻：可用 Pod ≥ replicas - maxUnavailable
- 任何时刻：总 Pod ≤ replicas + maxSurge

### 4. 服务发现

```
Service "web-service"
├── Selector: {app: web}
├── Port: 80 → TargetPort: 8080
└── Endpoints:
    ├── pod-abc (10.0.1.5:8080) ← Running + Ready
    ├── pod-def (10.0.1.8:8080) ← Running + Ready
    └── (pod-ghi 不在列表中)    ← Running 但 Not Ready

客户端 → "web-service:80" → 负载均衡 → pod-abc 或 pod-def
```

### 5. 自愈（Self-Healing）

```
正常状态: 3 个 Pod 运行中
  [●] [●] [●]

Pod 崩溃:
  [●] [✗] [●]

控制循环检测到 current(2) < desired(3):
  [●] [○] [●]  ← 自动创建新 Pod

新 Pod 就绪:
  [●] [●] [●]  ← 恢复！
```

### 6. Pod 调度（Filter → Score → Bind）

```
10 个节点

过滤阶段（Filter）:
  node-1: ✗ Offline
  node-2: ✗ 内存不足
  node-3: ✓
  node-4: ✓
  node-5: ✗ Degraded
  node-6: ✓
  ...
  候选: [node-3, node-4, node-6, node-7, node-9]

打分阶段（Score）:
  node-3: 资源均衡=40 + 拓扑分散=45 = 85
  node-4: 资源均衡=35 + 拓扑分散=50 = 85
  node-6: 资源均衡=50 + 拓扑分散=40 = 90  ← 最高分
  node-7: 资源均衡=30 + 拓扑分散=35 = 65
  node-9: 资源均衡=45 + 拓扑分散=30 = 75

绑定阶段（Bind）:
  Pod → node-6
```

## 练习任务

1. `pkg/orchestrator/types.go` — 阅读理解所有类型定义
2. `pkg/orchestrator/scheduler.go` — 实现 Pod 调度（Filter → Score → Bind）
3. `pkg/orchestrator/replicaset.go` — 实现控制循环（最核心！）
4. `pkg/orchestrator/deployment.go` — 实现滚动更新和回滚
5. `pkg/orchestrator/service.go` — 实现服务发现和端点管理

## 验证

```bash
go test ./pkg/orchestrator/...
go test -race ./pkg/orchestrator/...
```
