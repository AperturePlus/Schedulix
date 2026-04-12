package scheduler

import "schedulix/pkg/model"

// BestFitStrategy Best-Fit 调度算法。
//
// 算法描述：
//   遍历所有可用节点，选择满足资源需求且剩余资源最少的节点。
//   "最少剩余"以可用内存为主要指标。
//
// 特点：
//   - 时间复杂度：O(n)，需要遍历所有节点
//   - 优点：资源利用率高，减少碎片化
//   - 缺点：调度延迟略高于 First-Fit
//   - 适用场景：资源紧张，需要最大化利用率
//
// 与 First-Fit 的对比：
//   First-Fit 找到第一个就返回，Best-Fit 要看完所有才决定。
//   Best-Fit 倾向于把小任务塞进"刚好够"的节点，留大节点给大任务。
type BestFitStrategy struct{}

// Schedule 执行 Best-Fit 调度。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 调用 ValidateScheduleInput(task, cluster) — 不信任调用方
// 2. 获取 Idle 节点列表，如果为空则 fallback 到 Busy 节点
// 3. 遍历节点，找到满足 CanAccept 且 AvailableMemory() 最小的节点
//    - 注意：AvailableMemory() 已有防御性处理（不会返回负值）
// 4. 如果有多个节点剩余资源相同，选第一个（确定性行为，便于测试）
// 5. 返回该节点 ID
// 6. 如果没有满足条件的节点，返回 ErrNoAvailableNode
func (s *BestFitStrategy) Schedule(task *model.Task, cluster *model.Cluster) (string, error) {
	// TODO: 实现
	panic("not implemented")
}

func (s *BestFitStrategy) Name() string {
	return "best-fit"
}
