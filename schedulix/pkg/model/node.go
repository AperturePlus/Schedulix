package model

import (
	"errors"
	"fmt"
)

// ─── 错误定义 ───────────────────────────────────────────────
// 合格的程序员为每种失败模式定义明确的错误类型，
// 调用方可以用 errors.Is() 精确判断失败原因并做出不同的应对。

var (
	// ErrNilNode 节点指针为 nil。
	// 防御性编程：任何接收 *GPU_Node 的方法都应检查 nil。
	ErrNilNode = errors.New("gpu node is nil")

	// ErrInvalidNodeID 节点 ID 为空或不合法。
	ErrInvalidNodeID = errors.New("invalid node ID: must be non-empty")

	// ErrNegativeResource 资源值为负数（数据损坏或 bug 的信号）。
	ErrNegativeResource = errors.New("resource value must not be negative")

	// ErrFaultRateOutOfRange 故障率不在 [0.0, 1.0] 范围内。
	ErrFaultRateOutOfRange = errors.New("fault rate must be in [0.0, 1.0]")

	// ErrMemoryOvercommit 已用内存超过总内存（不变量被破坏）。
	ErrMemoryOvercommit = errors.New("memory used exceeds memory total")
)

// GPU_Node 模拟的 GPU 计算节点。
//
// 字段说明：
//   - ComputePower, MemoryTotal: 静态属性（创建后不变）
//   - MemoryUsed, AssignedTasks: 动态属性（随任务分配变化）
//   - FaultRate: 用于事件模拟器的伯努利试验，取值 [0.0, 1.0]
//   - RackID/CabinetID/DataCenterID: 拓扑信息，冗余存储避免反向查找
//
// 不变量（Invariants）— 任何时刻都必须成立：
//   - ID 非空
//   - MemoryUsed <= MemoryTotal
//   - MemoryUsed >= 0, MemoryTotal >= 0, ComputePower >= 0
//   - FaultRate ∈ [0.0, 1.0]
//   - len(AssignedTasks) 与实际分配的任务数一致
type GPU_Node struct {
	ID            string     `json:"id"`
	Status        NodeStatus `json:"status"`
	ComputePower  int        `json:"compute_power"`  // 算力（TFLOPS）
	MemoryTotal   int        `json:"memory_total"`   // 总内存（MB）
	MemoryUsed    int        `json:"memory_used"`    // 已用内存（MB）
	FaultRate     float64    `json:"fault_rate"`      // 故障率 [0.0, 1.0]
	AssignedTasks []string   `json:"assigned_tasks"`  // 已分配任务 ID 列表
	FaultCount    int        `json:"fault_count"`     // 累计故障次数
	UptimeMs      int64      `json:"uptime_ms"`       // 累计运行时间（毫秒）
	RackID        string     `json:"rack_id"`         // 所属机架 ID
	CabinetID     string     `json:"cabinet_id"`      // 所属机柜 ID
	DataCenterID  string     `json:"data_center_id"`  // 所属数据中心 ID
}

// Validate 检查节点数据的完整性和一致性。
//
// TODO(learner): 实现此方法
// 鲁棒性要求 — 检查所有不变量：
// 1. ID 非空 → 否则返回 ErrInvalidNodeID
// 2. ComputePower >= 0, MemoryTotal >= 0, MemoryUsed >= 0 → 否则返回 ErrNegativeResource
// 3. MemoryUsed <= MemoryTotal → 否则返回 ErrMemoryOvercommit
// 4. FaultRate ∈ [0.0, 1.0] → 否则返回 ErrFaultRateOutOfRange
//
// 提示：使用 fmt.Errorf("%w: detail", ErrXxx) 包装错误，保留上下文信息
func (n *GPU_Node) Validate() error {
	// TODO: 实现
	panic("not implemented")
}

// AvailableMemory 返回节点剩余可用内存（MB）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 如果 MemoryUsed > MemoryTotal（不变量被破坏），返回 0 而非负数
//    — 防御性编程：即使内部状态异常，也不向外暴露负值
// 2. 这是一个"安全降级"的例子：检测到异常但不 panic，返回安全的默认值
func (n *GPU_Node) AvailableMemory() int {
	// TODO: 实现
	panic("not implemented")
}

// CanAccept 判断节点是否能接受指定资源需求的任务。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 如果节点状态为 Offline，直接返回 false
// 2. 如果资源需求为负值（调用方 bug），返回 false — 不信任输入
// 3. 计算有效算力：Degraded 状态下算力减半
// 4. 检查有效算力 >= 需求算力 且 可用内存 >= 需求内存
// 5. 使用 AvailableMemory()（已有防御性处理）而非直接计算
func (n *GPU_Node) CanAccept(req ResourceRequirement) bool {
	// TODO: 实现
	panic("not implemented")
}

// AllocateTask 将任务分配到此节点，更新资源占用。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 先调用 CanAccept 检查 — 不信任调用方已经检查过
// 2. 检查 taskID 非空
// 3. 检查 taskID 是否已在 AssignedTasks 中（幂等性：重复分配不应重复扣资源）
// 4. 更新 MemoryUsed += req.Memory
// 5. 追加 taskID 到 AssignedTasks
// 6. 返回 error 而非 panic — 让调用方决定如何处理
func (n *GPU_Node) AllocateTask(taskID string, req ResourceRequirement) error {
	// TODO: 实现
	panic("not implemented")
}

// ReleaseTask 释放任务占用的资源。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 检查 taskID 是否在 AssignedTasks 中 — 释放不存在的任务应返回错误而非静默忽略
// 2. MemoryUsed -= req.Memory，但结果不能为负（clamp to 0）
// 3. 从 AssignedTasks 中移除 taskID
// 4. 如果 AssignedTasks 为空且状态为 Busy，考虑将状态改为 Idle
func (n *GPU_Node) ReleaseTask(taskID string, req ResourceRequirement) error {
	// TODO: 实现
	panic("not implemented")
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
var _ = errors.New
