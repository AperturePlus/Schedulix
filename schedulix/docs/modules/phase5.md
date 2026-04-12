# 阶段五：负载均衡与万卡规模

## 学习目标

实现负载均衡策略，优化系统支持万卡（10,000 节点）规模，掌握性能优化技术。

## 前置知识

- 阶段四完成
- 了解标准差的概念

## 核心概念

### 1. 静态负载均衡 — 加权随机选择

按节点算力权重分配任务。算力越高，被选中概率越大。

**算法：加权随机选择（Weighted Random Selection）**

```
节点:  A(算力100)  B(算力200)  C(算力50)
总算力: 350

生成随机数 r ∈ [0, 350)

r=0~99    → 选 A (概率 100/350 ≈ 28.6%)
r=100~299 → 选 B (概率 200/350 ≈ 57.1%)
r=300~349 → 选 C (概率 50/350  ≈ 14.3%)
```

```go
func weightedSelect(nodes []*GPU_Node) *GPU_Node {
    totalPower := 0
    for _, n := range nodes {
        totalPower += n.ComputePower
    }
    
    r := rand.Intn(totalPower)
    cumulative := 0
    for _, n := range nodes {
        cumulative += n.ComputePower
        if r < cumulative {
            return n
        }
    }
    return nodes[len(nodes)-1] // 兜底
}
```

### 2. 动态负载均衡 — 最低负载优先

选择当前负载最低的节点。

**负载计算**：
```
load = MemoryUsed / MemoryTotal
```

**重平衡判断 — 标准差（Standard Deviation）**：

标准差衡量数据的离散程度。节点负载的标准差越大，说明负载越不均匀。

```
节点负载: [0.2, 0.8, 0.3, 0.7]

1. 计算均值: mean = (0.2+0.8+0.3+0.7) / 4 = 0.5
2. 计算方差: variance = ((0.2-0.5)² + (0.8-0.5)² + (0.3-0.5)² + (0.7-0.5)²) / 4
                      = (0.09 + 0.09 + 0.04 + 0.04) / 4
                      = 0.065
3. 标准差:   stddev = √0.065 ≈ 0.255

如果 threshold = 0.2，则 0.255 > 0.2 → 需要重平衡
```

```go
import "math"

func stddev(loads []float64) float64 {
    n := float64(len(loads))
    mean := 0.0
    for _, l := range loads {
        mean += l
    }
    mean /= n
    
    variance := 0.0
    for _, l := range loads {
        diff := l - mean
        variance += diff * diff
    }
    variance /= n
    
    return math.Sqrt(variance)
}
```

### 3. 万卡规模性能优化

10,000 个节点时，朴素遍历可能太慢。关键优化：

#### 辅助索引

```go
// 不优化：每次调度遍历 10,000 个节点
for _, node := range cluster.Nodes {
    if node.Status == NodeStatusIdle { ... }
}

// 优化：只遍历 Idle 节点（通常远少于 10,000）
idleNodeIDs := cluster.statusIndex[NodeStatusIdle]
for _, id := range idleNodeIDs {
    node := cluster.Nodes[id]
    ...
}
```

#### Benchmark 测试

Go 内置 benchmark 支持：

```go
func BenchmarkSchedule10000(b *testing.B) {
    cluster := model.NewCluster(10000)
    strategy := &FirstFitStrategy{}
    task := &model.Task{Resource: model.ResourceRequirement{Memory: 1024}}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        strategy.Schedule(task, cluster)
    }
}
```

```bash
go test -bench=BenchmarkSchedule10000 ./pkg/scheduler/...
# 输出示例：
# BenchmarkSchedule10000-8    50000    25000 ns/op
# 意思：每次调度约 25μs（远小于 100ms 要求）
```

### 4. 集群拓扑

```
数据中心 (DataCenter)
  └── 机柜 (Cabinet)
        └── 机架 (Rack)
              └── 节点 (GPU_Node)

示例：2 个数据中心 × 5 个机柜 × 10 个机架 × 100 个节点 = 10,000 节点
```

拓扑信息可用于**拓扑感知调度**：优先将关联任务调度到同一机架（减少网络延迟）。

## 练习任务

1. 打开 `pkg/balancer/static.go`，实现加权随机选择
2. 打开 `pkg/balancer/dynamic.go`，实现最低负载选择和标准差重平衡判断
3. 扩展 `pkg/model/cluster.go` 的 `NewCluster`，支持万卡规模和拓扑结构
4. 实现 `SnapshotToJSON` / `RestoreFromJSON`
5. 编写 benchmark 测试，验证调度延迟 < 100ms

## 验证

```bash
go test ./pkg/balancer/...
go test -bench=. ./pkg/model/...
go test -bench=. ./pkg/scheduler/...
```
