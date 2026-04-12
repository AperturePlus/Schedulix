# 阶段四：事件模拟与容灾恢复

## 学习目标

构建事件模拟引擎和容灾恢复机制，掌握事件驱动编程、概率模型和观察者模式。

## 前置知识

- 阶段三完成
- 了解概率的基本概念

## 核心概念

### 1. 离散事件模拟（Discrete Event Simulation, DES）

DES 是模拟真实系统的经典方法。核心思想：**时间不是均匀流逝的，而是跳到下一个事件发生的时刻**。

```
传统时间步进（每步都检查）：
t=0  t=1  t=2  t=3  t=4  t=5  t=6  t=7  t=8  t=9
 ✓    ✓    ✓    ✓    ✓    ✓    ✓    ✓    ✓    ✓
                 ↑              ↑
               事件1          事件2

离散事件模拟（只处理有事件的时刻）：
t=0 ──────→ t=3 ──────→ t=6 ──────→ ...
            事件1        事件2
            （跳过空闲时段）
```

**Schedulix 同时支持两种模式**：
- 时间步进模式（入门）：简单直观，逐步推进
- 离散事件队列模式（进阶）：高效，适合万卡规模

### 2. 伯努利试验（Bernoulli Trial）

每次试验只有两个结果：发生 / 不发生，以概率 p 决定。

```go
// 对每个节点，每个时间步，独立判定是否故障
if rand.Float64() < config.NodeDownProb {
    // 触发宕机事件！
}
```

**为什么各事件概率独立？**
一个节点可能同时遭遇网络延迟和性能降级（虽然概率很低）。独立概率模型更贴近现实。

**概率参数的直觉**：
- `NodeDownProb = 0.005`：每步每节点有 0.5% 概率宕机
- 10,000 节点 × 0.005 = 平均每步约 50 个节点宕机
- 这模拟了大规模集群中"总有节点在出问题"的现实

### 3. 观察者模式（Observer Pattern）

事件模拟器不直接处理故障，而是通知注册的处理器：

```
EventSimulator
    │
    ├──→ RecoveryEngine.OnFault()    // 处理任务迁移
    ├──→ MetricsCollector.OnFault()  // 记录指标
    └──→ Logger.OnFault()            // 记录日志
```

```go
type EventHandler interface {
    OnFault(event *FaultEvent) error
    OnRecovery(event *FaultEvent) error
}

// 注册
simulator.RegisterHandler(recoveryEngine)
simulator.RegisterHandler(metricsCollector)

// 事件发生时，遍历通知
for _, handler := range es.handlers {
    handler.OnFault(event)
}
```

### 4. 检查点机制（Checkpoint）

定期保存任务执行进度，故障时从最近检查点恢复，而非从头开始。

```
任务执行进度：0% ──→ 25% ──→ 50% ──→ 75% ──→ 100%
                      ↑         ↑
                   检查点1    检查点2

节点在 60% 时宕机：
  无检查点：从 0% 重新开始 😢
  有检查点：从 50% 恢复    😊
```

### 5. 容灾恢复流程

```
节点宕机
  │
  ▼
检测受影响任务（该节点的 AssignedTasks）
  │
  ▼
对每个任务：
  ├── 查询最近检查点 → 恢复进度
  ├── MigrationCount++
  ├── MigrationCount >= 3？
  │     ├── 是 → 标记 Failed，发出告警
  │     └── 否 → 重新入队，等待重新调度
  │
  ▼
记录恢复日志
```

## 练习任务

1. 打开 `pkg/simulator/config.go`，实现 `Validate` 和 `ParseConfig`
2. 打开 `pkg/simulator/engine.go`，实现事件模拟引擎（先实现时间步进模式）
3. 打开 `pkg/recovery/checkpoint.go`，实现检查点的 Save/Load/Delete
4. 打开 `pkg/recovery/engine.go`，实现 `OnFault` 和 `OnRecovery`

## 验证

```bash
go test ./pkg/simulator/...
go test ./pkg/recovery/...
```
