package balancer

import (
	"errors"

	"schedulix/pkg/model"
)

// ErrNoSuitableNode 没有合适的节点。
var ErrNoSuitableNode = errors.New("no suitable node for load balancing")

// BalanceStrategy 负载均衡策略接口。
//
// 与 ScheduleStrategy 的区别：
//   - ScheduleStrategy: "新任务放哪里"
//   - BalanceStrategy: "已有任务是否需要重新分布"
type BalanceStrategy interface {
	// SelectNode 根据负载情况选择最优节点。
	// nodes 已经过状态过滤（仅 Idle 或 Busy）。
	SelectNode(task *model.Task, nodes []*model.GPU_Node) (nodeID string, err error)

	// ShouldRebalance 判断是否需要触发重平衡。
	// threshold 取值 [0.0, 1.0]，0.0 = 任何不均衡都触发，1.0 = 从不触发。
	ShouldRebalance(nodes []*model.GPU_Node, threshold float64) bool

	// Name 返回策略名称。
	Name() string
}
