package balancer

import "schedulix/pkg/model"

// StaticBalancer 静态负载均衡策略。
//
// 算法描述：
//   基于节点算力权重分配任务。算力越高的节点，被选中的概率越大。
//
// 实现方式（加权随机选择 Weighted Random Selection）：
//   1. 计算所有候选节点的算力总和 totalPower
//   2. 生成 [0, totalPower) 的随机数 r
//   3. 遍历节点，累加算力，当累加值 > r 时选中该节点
type StaticBalancer struct{}

// SelectNode 基于算力权重选择节点。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. nodes == nil 或 len(nodes) == 0 → 返回 ErrNoSuitableNode
// 2. 过滤掉 nil 节点和不满足 CanAccept 的节点
// 3. 过滤后为空 → 返回 ErrNoSuitableNode
// 4. 计算 totalPower — 如果所有节点算力为 0（异常数据），fallback 到均匀随机选择
// 5. 加权随机选择
// 6. 兜底：如果循环结束没选中（浮点精度问题），返回最后一个候选节点
func (s *StaticBalancer) SelectNode(task *model.Task, nodes []*model.GPU_Node) (string, error) {
	// TODO: 实现加权随机选择
	panic("not implemented")
}

// ShouldRebalance 静态策略不主动触发重平衡。
func (s *StaticBalancer) ShouldRebalance(nodes []*model.GPU_Node, threshold float64) bool {
	return false
}

func (s *StaticBalancer) Name() string {
	return "static"
}
