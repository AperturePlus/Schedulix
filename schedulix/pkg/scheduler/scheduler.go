package scheduler

import (
	"errors"
	"fmt"
	"sync"

	"schedulix/pkg/model"
	"schedulix/pkg/queue"
)

// ─── 错误定义 ───────────────────────────────────────────────

var (
	// ErrSchedulerNotReady 调度器未正确初始化。
	ErrSchedulerNotReady = errors.New("scheduler is not ready: missing strategy, queue, or cluster")

	// ErrScheduleFailed 调度失败（包装底层错误）。
	ErrScheduleFailed = errors.New("schedule failed")

	// ErrResourceUpdateFailed 资源更新失败（调度成功但资源分配失败 — 需要回滚）。
	ErrResourceUpdateFailed = errors.New("resource update failed after scheduling")
)

// Scheduler 调度器核心，持有策略、队列和集群的引用。
//
// 鲁棒性设计：
//   - 所有依赖在使用前检查 nil
//   - 调度失败时任务不丢失（重新入队）
//   - 资源分配失败时回滚（释放已预留资源）
//   - 支持运行时切换策略（加锁保护）
type Scheduler struct {
	strategy ScheduleStrategy
	queue    *queue.TaskQueue
	cluster  *model.Cluster
	mu       sync.RWMutex // 保护 strategy 的运行时切换
}

// NewScheduler 创建调度器实例。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - strategy, q, c 任一为 nil → 仍然创建实例但标记为未就绪
//   （延迟检查：在 ScheduleNext 时返回 ErrSchedulerNotReady）
// - 这样做的好处：允许分步初始化，不强制所有依赖同时就绪
func NewScheduler(strategy ScheduleStrategy, q *queue.TaskQueue, c *model.Cluster) *Scheduler {
	return &Scheduler{
		strategy: strategy,
		queue:    q,
		cluster:  c,
	}
}

// SetStrategy 运行时切换调度策略（线程安全）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - strategy == nil → 返回 ErrNilStrategy（不允许设置为 nil）
// - 使用 s.mu.Lock() 保护写操作
func (s *Scheduler) SetStrategy(strategy ScheduleStrategy) error {
	// TODO: 实现
	panic("not implemented")
}

// isReady 检查调度器是否就绪。
func (s *Scheduler) isReady() error {
	if s.strategy == nil || s.queue == nil || s.cluster == nil {
		return ErrSchedulerNotReady
	}
	return nil
}

// ScheduleNext 从队列取出最高优先级任务并调度。
//
// TODO(learner): 实现此方法
// 鲁棒性要求（完整的错误处理链）：
// 1. 检查 isReady() → 不就绪则返回错误
// 2. 从 queue 中 Dequeue 一个任务
//    - 队列为空 → 返回 nil, queue.ErrQueueEmpty（正常情况，不是错误）
// 3. 调用 strategy.Schedule(task, cluster) 选择节点
//    - 失败（ErrNoAvailableNode）→ 将任务重新 Enqueue（不丢失任务！）
//    - Enqueue 也失败 → 返回任务和包装错误（调用方可以决定如何处理这个任务）
// 4. 成功 → 调用 node.AllocateTask(task.ID, task.Resource) 更新资源
//    - AllocateTask 失败 → 回滚：不更新任务状态，重新入队
// 5. 更新任务状态为 Running，设置 AssignedNodeID
// 6. 返回任务和 nil
//
// 关键原则：任务绝不能丢失。任何失败路径都必须确保任务回到队列或返回给调用方。
func (s *Scheduler) ScheduleNext() (*model.Task, error) {
	// TODO: 实现
	panic("not implemented")
}

// ScheduleTask 调度指定任务（不从队列取，直接调度）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 检查 isReady()
// 2. task == nil → 返回 "", ErrNilTask
// 3. 调用 strategy.Schedule(task, cluster)
// 4. 成功 → 调用 node.AllocateTask，失败则回滚
// 5. 更新任务状态
// 6. 返回 nodeID, nil
func (s *Scheduler) ScheduleTask(task *model.Task) (string, error) {
	// TODO: 实现
	panic("not implemented")
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
