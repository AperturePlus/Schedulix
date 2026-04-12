package metrics

import (
	"sync"

	"schedulix/pkg/model"
)

// MetricsCollector 指标采集器。
//
// 鲁棒性设计：
//   - 环形缓冲区固定大小，不会 OOM
//   - bufferSize 不合法时使用默认值
//   - cluster 为 nil 时 Collect 静默返回（不 panic）
//   - 所有公开方法线程安全
//   - GetHistory 返回副本
type MetricsCollector struct {
	cluster *model.Cluster

	// 环形缓冲区
	clusterSnapshots []model.ClusterMetrics
	nodeSnapshots    [][]model.NodeMetrics
	bufferSize       int
	cursor           int
	count            int // 已写入总数（用于判断缓冲区是否已满）
	version          int

	mu sync.RWMutex
}

// DefaultBufferSize 默认缓冲区大小。
const DefaultBufferSize = 1000

// NewMetricsCollector 创建指标采集器。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - cluster == nil → 仍然创建，Collect 时检查
// - bufferSize <= 0 → 使用 DefaultBufferSize
// - 初始化环形缓冲区（make slice with bufferSize length）
func NewMetricsCollector(cluster *model.Cluster, bufferSize int) *MetricsCollector {
	// TODO: 实现
	panic("not implemented")
}

// Collect 执行一次指标采集。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. cluster == nil → 静默返回（不 panic，不采集）
// 2. 使用 mu.Lock() 保护写操作
// 3. 遍历集群节点，计算集群级指标和节点级指标
//    - 跳过 nil 节点
//    - 除法运算检查除数为 0（totalMemory == 0 → utilization = 0）
// 4. 递增 version
// 5. 写入环形缓冲区的 cursor 位置
// 6. cursor = (cursor + 1) % bufferSize
// 7. count++（但不超过 bufferSize，用于 GetHistory 判断有效数据量）
func (mc *MetricsCollector) Collect() {
	// TODO: 实现
	panic("not implemented")
}

// GetLatestClusterMetrics 获取最新的集群指标快照。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 使用 mu.RLock()
// 2. count == 0（从未采集过）→ 返回 nil
// 3. 返回副本（值类型，直接赋值即可）
func (mc *MetricsCollector) GetLatestClusterMetrics() *model.ClusterMetrics {
	// TODO: 实现
	panic("not implemented")
}

// GetHistory 获取历史指标（最近 n 条）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. n <= 0 → 返回空切片
// 2. n > count → n = count（不返回未初始化的数据）
// 3. n > bufferSize → n = bufferSize
// 4. 从环形缓冲区中按时间顺序提取最近 n 条
// 5. 返回副本
func (mc *MetricsCollector) GetHistory(n int) []model.ClusterMetrics {
	// TODO: 实现
	panic("not implemented")
}
