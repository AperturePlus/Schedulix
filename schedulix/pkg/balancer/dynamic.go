package balancer

import "schedulix/pkg/model"

// DynamicBalancer 动态负载均衡策略。
//
// 算法描述：
//   基于节点实时资源使用率分配任务。选择当前负载最低的节点。
//   当节点间负载差异超过阈值时，触发任务迁移（重平衡）。
type DynamicBalancer struct{}

// SelectNode 选择当前负载最低的节点。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. nodes == nil 或 len(nodes) == 0 → 返回 ErrNoSuitableNode
// 2. 过滤掉 nil 节点和不满足 CanAccept 的节点
// 3. 过滤后为空 → 返回 ErrNoSuitableNode
// 4. 计算每个节点的负载 = MemoryUsed / MemoryTotal
//    - MemoryTotal == 0（异常数据）→ 视为负载 1.0（满载），跳过该节点
//    - 使用 float64 除法避免整数除法截断
// 5. 返回负载最低的节点 ID
// 6. 多个节点负载相同 → 选第一个（确定性行为）
func (d *DynamicBalancer) SelectNode(task *model.Task, nodes []*model.GPU_Node) (string, error) {
	// TODO: 实现
	panic("not implemented")
}

// ShouldRebalance 判断是否需要重平衡。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. nodes == nil 或 len(nodes) <= 1 → 返回 false（无法计算标准差）
// 2. threshold < 0 → 视为 0（总是触发）
// 3. threshold > 1 → 视为 1（从不触发）
// 4. 计算节点负载的标准差
//    - MemoryTotal == 0 的节点 → 负载视为 1.0
// 5. return stddev > threshold
//
// 标准差公式：
//   mean = sum(load_i) / n
//   variance = sum((load_i - mean)^2) / n
//   stddev = sqrt(variance)
func (d *DynamicBalancer) ShouldRebalance(nodes []*model.GPU_Node, threshold float64) bool {
	// TODO: 实现
	panic("not implemented")
}

func (d *DynamicBalancer) Name() string {
	return "dynamic"
}
