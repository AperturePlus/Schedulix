package model

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// ─── 错误定义 ───────────────────────────────────────────────

var (
	// ErrClusterNil 集群指针为 nil。
	ErrClusterNil = errors.New("cluster is nil")

	// ErrNodeNotFound 节点 ID 在集群中不存在。
	ErrNodeNotFound = errors.New("node not found in cluster")

	// ErrInvalidClusterSize 集群大小不合法。
	ErrInvalidClusterSize = errors.New("cluster size must be positive")

	// ErrDuplicateNodeID 节点 ID 重复。
	ErrDuplicateNodeID = errors.New("duplicate node ID in cluster")

	// ErrCorruptedIndex 辅助索引与主数据不一致（内部错误）。
	ErrCorruptedIndex = errors.New("status index is corrupted: node exists in index but not in nodes map")

	// ErrInvalidSnapshot JSON 快照数据不合法或已损坏。
	ErrInvalidSnapshot = errors.New("invalid or corrupted cluster snapshot")
)

// DataCenter 数据中心 — 拓扑最顶层
type DataCenter struct {
	ID       string    `json:"id"`
	Cabinets []Cabinet `json:"cabinets"`
}

// Cabinet 机柜 — 拓扑中间层
type Cabinet struct {
	ID    string `json:"id"`
	Racks []Rack `json:"racks"`
}

// Rack 机架 — 拓扑最底层，直接包含节点引用
type Rack struct {
	ID      string   `json:"id"`
	NodeIDs []string `json:"node_ids"`
}

// Cluster 模拟集群。
//
// 并发安全策略：
//   - mu (sync.RWMutex) 保护 Nodes 和 statusIndex 的并发访问
//   - 读操作使用 RLock，写操作使用 Lock
//   - mu 和 statusIndex 不参与 JSON 序列化
//
// 不变量（Invariants）：
//   - Nodes map 中的每个节点 ID 必须与节点的 ID 字段一致
//   - statusIndex 必须与 Nodes 中各节点的 Status 保持同步
//   - 任何修改 Nodes 或节点 Status 的操作都必须同步更新 statusIndex
type Cluster struct {
	Nodes       map[string]*GPU_Node    `json:"nodes"`
	DataCenters []DataCenter            `json:"data_centers"`
	statusIndex map[NodeStatus][]string // 按状态分组的节点 ID 索引
	mu          sync.RWMutex
}

// NewCluster 创建包含 count 个节点的模拟集群。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. count <= 0 → 返回 nil 并记录错误（或返回空集群，你来决定策略）
// 2. 创建 count 个 GPU_Node，每个节点具有唯一 ID（如 "node-0001"）
// 3. 初始状态为 Idle
// 4. 初始化 Nodes map（预分配容量 count，避免 rehash）和 statusIndex
// 5. 确保 statusIndex 与 Nodes 一致
//
// 提示：使用 fmt.Sprintf("node-%04d", i) 生成唯一 ID
func NewCluster(count int) *Cluster {
	// TODO: 实现
	panic("not implemented")
}

// GetNode 根据 ID 获取节点（线程安全）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 使用 c.mu.RLock() / defer c.mu.RUnlock() 保护读操作
// 2. nodeID 为空 → 返回 nil, false（不 panic）
// 3. 节点不存在 → 返回 nil, false
func (c *Cluster) GetNode(nodeID string) (*GPU_Node, bool) {
	// TODO: 实现
	panic("not implemented")
}

// GetAvailableNodes 返回指定状态的节点列表（线程安全）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 使用 c.mu.RLock() 保护读操作
// 2. 从 statusIndex 获取该状态的节点 ID 列表
// 3. 对每个 ID，从 Nodes map 中获取节点指针
// 4. 如果 ID 在 statusIndex 中但不在 Nodes map 中（索引损坏），
//    跳过该节点并记录警告（不 panic，优雅降级）
// 5. 返回空切片而非 nil（调用方可以安全 range）
func (c *Cluster) GetAvailableNodes(status NodeStatus) []*GPU_Node {
	// TODO: 实现
	panic("not implemented")
}

// UpdateNodeStatus 更新节点状态并同步更新辅助索引（线程安全）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 使用 c.mu.Lock() / defer c.mu.Unlock() 保护写操作
// 2. nodeID 为空 → 返回 fmt.Errorf("%w: empty node ID", ErrNodeNotFound)
// 3. 节点不存在 → 返回 ErrNodeNotFound（包装节点 ID 信息）
// 4. 新旧状态相同 → 直接返回 nil（幂等操作）
// 5. 从旧状态的索引中移除该节点 ID
// 6. 添加到新状态的索引中
// 7. 更新节点的 Status 字段
// 8. 如果旧状态索引中找不到该节点（索引不一致），仍然完成更新但记录警告
func (c *Cluster) UpdateNodeStatus(nodeID string, newStatus NodeStatus) error {
	// TODO: 实现
	panic("not implemented")
}

// FilterByStatus 返回指定状态的节点切片（不加锁，调用方需持有锁）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 遍历 c.Nodes，筛选 Status == status 的节点
// 2. 跳过 nil 节点（防御性：map 中不应有 nil value，但以防万一）
// 3. 返回空切片而非 nil
func (c *Cluster) FilterByStatus(status NodeStatus) []*GPU_Node {
	// TODO: 实现
	panic("not implemented")
}

// SortByComputePower 将节点切片按算力降序排序。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. nodes 为 nil 或空 → 直接返回（不 panic）
// 2. 使用 sort.Slice，比较函数为 nodes[i].ComputePower > nodes[j].ComputePower
func SortByComputePower(nodes []*GPU_Node) {
	if len(nodes) == 0 {
		return
	}
	// TODO: 实现
	sort.Slice(nodes, func(i, j int) bool {
		panic("not implemented")
	})
}

// SnapshotToJSON 将集群状态序列化为 JSON 字节。
//
// TODO(learner): 实现此方法（阶段一进阶）
// 鲁棒性要求：
// 1. 使用 c.mu.RLock() / defer c.mu.RUnlock() 保护读操作
// 2. 使用 encoding/json.Marshal 序列化
// 3. 序列化失败 → 包装错误返回，不 panic
func (c *Cluster) SnapshotToJSON() ([]byte, error) {
	// TODO: 实现
	panic("not implemented")
}

// RestoreFromJSON 从 JSON 字节恢复集群状态。
//
// TODO(learner): 实现此方法（阶段一进阶）
// 鲁棒性要求：
// 1. data 为 nil 或空 → 返回 ErrInvalidSnapshot
// 2. 使用 encoding/json.Unmarshal 反序列化
// 3. 反序列化失败 → 包装错误返回
// 4. 验证恢复后的数据：
//    a. Nodes map 不为空
//    b. 每个节点的 ID 与 map key 一致
//    c. 每个节点通过 Validate()
// 5. 重建 statusIndex（遍历所有节点，按状态分组）
// 6. 初始化 mu（sync.RWMutex 零值即可用）
// 7. 任何验证失败 → 返回 ErrInvalidSnapshot + 具体原因
func RestoreFromJSON(data []byte) (*Cluster, error) {
	// TODO: 实现
	panic("not implemented")
}

// rebuildStatusIndex 重建辅助索引。
// 在 RestoreFromJSON 和索引损坏修复时使用。
//
// TODO(learner): 实现此方法
// 提示：遍历 Nodes，按 Status 分组收集节点 ID
func (c *Cluster) rebuildStatusIndex() {
	// TODO: 实现
	panic("not implemented")
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
