package scheduler

import "schedulix/pkg/model"

// RoundRobinStrategy Round-Robin 调度算法。
//
// 算法描述：
//   维护一个游标（cursor），每次调度从游标位置开始，
//   轮询可用节点，找到第一个满足资源需求的节点后分配，
//   游标移动到下一个位置。
//
// 特点：
//   - 时间复杂度：O(1) 均摊（大多数情况下很快找到可用节点）
//   - 优点：负载均匀分布，公平性好
//   - 缺点：不考虑节点剩余资源差异，可能浪费大节点
//   - 适用场景：节点配置相近，追求公平性
//
// 注意：cursor 不是线程安全的，并发场景需要在 ConcurrentScheduler 中加锁。
type RoundRobinStrategy struct {
	cursor int // 当前轮询位置
}

// Schedule 执行 Round-Robin 调度。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 调用 ValidateScheduleInput(task, cluster) — 不信任调用方
// 2. 获取 Idle 节点列表，如果为空则 fallback 到 Busy 节点
// 3. 如果合并后仍为空，返回 ErrNoAvailableNode
// 4. cursor 可能大于当前节点数（节点数量动态变化），用 cursor % len(nodes) 归位
// 5. 从归位后的位置开始，遍历一圈
// 6. 找到第一个满足 CanAccept 的节点，更新 cursor = (pos + 1)，返回节点 ID
// 7. 遍历一圈都没找到，返回 ErrNoAvailableNode
//
// 防御性要点：
// - 节点列表在两次调度之间可能变化（节点宕机/恢复），cursor 必须取模
// - 不要假设 cursor 的值在合法范围内
func (s *RoundRobinStrategy) Schedule(task *model.Task, cluster *model.Cluster) (string, error) {
	// TODO: 实现
	panic("not implemented")
}

func (s *RoundRobinStrategy) Name() string {
	return "round-robin"
}
