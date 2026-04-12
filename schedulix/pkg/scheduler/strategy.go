package scheduler

import (
	"errors"

	"schedulix/pkg/model"
)

// ─── 错误定义 ───────────────────────────────────────────────
// 调度器的每种失败模式都有明确的错误类型。
// 调用方可以用 errors.Is() 区分"没有可用节点"和"输入不合法"，
// 从而做出不同的应对（等待重试 vs 拒绝请求）。

var (
	// ErrNoAvailableNode 没有可用节点满足资源需求。
	// 这是一个可恢复错误：稍后重试可能成功（有节点释放资源）。
	ErrNoAvailableNode = errors.New("no available node satisfies the resource requirement")

	// ErrNilTask 任务指针为 nil。
	ErrNilTask = errors.New("task is nil")

	// ErrNilCluster 集群指针为 nil。
	ErrNilCluster = errors.New("cluster is nil")

	// ErrNilStrategy 调度策略为 nil。
	ErrNilStrategy = errors.New("schedule strategy is nil")

	// ErrTaskAlreadyAssigned 任务已经被分配到节点（重复调度）。
	// 幂等性检查：防止同一任务被调度两次。
	ErrTaskAlreadyAssigned = errors.New("task is already assigned to a node")

	// ErrInvalidResourceRequest 资源需求不合法（负值或零值）。
	ErrInvalidResourceRequest = errors.New("invalid resource request: values must be positive")
)

// ScheduleStrategy 定义调度算法的统一抽象。
//
// 设计要点（策略模式 Strategy Pattern）：
//   - 每种调度算法实现此接口
//   - Scheduler 持有接口引用，运行时可切换策略
//   - Schedule 方法必须是无副作用的：不修改 task 或 cluster 状态
//   - 返回 nodeID 而非 *GPU_Node，资源更新由 Scheduler 统一执行
//
// 鲁棒性契约：
//   - 实现者必须检查 task 和 cluster 是否为 nil
//   - 实现者必须检查 task.Resource 是否合法
//   - 实现者在找不到节点时返回 ErrNoAvailableNode，不 panic
type ScheduleStrategy interface {
	// Schedule 将任务分配到合适的节点，返回目标节点 ID。
	// 如果没有可用节点，返回 ErrNoAvailableNode。
	Schedule(task *model.Task, cluster *model.Cluster) (nodeID string, err error)

	// Name 返回策略名称，用于日志和监控标识。
	Name() string
}

// ValidateScheduleInput 校验调度输入参数。
// 所有 ScheduleStrategy 实现应在 Schedule 方法开头调用此函数。
//
// TODO(learner): 实现此函数
// 检查：
// 1. task == nil → ErrNilTask
// 2. cluster == nil → ErrNilCluster
// 3. task.Resource.ComputePower <= 0 或 task.Resource.Memory <= 0 → ErrInvalidResourceRequest
// 4. task.AssignedNodeID != "" → ErrTaskAlreadyAssigned（已分配的任务不应再调度）
func ValidateScheduleInput(task *model.Task, cluster *model.Cluster) error {
	// TODO: 实现
	panic("not implemented")
}
