package scheduler

import "schedulix/pkg/model"

// FirstFitStrategy First-Fit 调度算法。
//
// 算法描述：
//   遍历所有可用节点（状态为 Idle 或 Busy），
//   返回第一个满足任务资源需求的节点。
//
// 特点：
//   - 时间复杂度：O(n) 最坏情况
//   - 优点：实现简单，调度速度快
//   - 缺点：可能导致资源碎片化，前面的节点负载偏高
//   - 适用场景：对调度延迟敏感，节点数量不大
type FirstFitStrategy struct{}

// Schedule 执行 First-Fit 调度。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 调用 ValidateScheduleInput(task, cluster) — 不信任调用方
// 2. 获取所有 Idle 状态的节点列表（cluster.GetAvailableNodes）
// 3. 如果列表为空 → 尝试 Busy 状态的节点作为 fallback（Busy 节点可能还有剩余资源）
// 4. 遍历节点，调用 node.CanAccept(task.Resource)
// 5. 返回第一个满足条件的节点 ID
// 6. 如果没有满足条件的节点，返回 ErrNoAvailableNode
func (s *FirstFitStrategy) Schedule(task *model.Task, cluster *model.Cluster) (string, error) {
	// TODO: 实现
	panic("not implemented")
}

func (s *FirstFitStrategy) Name() string {
	return "first-fit"
}
