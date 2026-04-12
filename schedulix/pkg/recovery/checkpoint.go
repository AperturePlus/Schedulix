package recovery

import (
	"sync"

	"schedulix/pkg/model"
)

// CheckpointStore 检查点存储。
// 使用内存 map，每个任务只保留最新一个检查点。
//
// 鲁棒性设计：
//   - 线程安全（sync.RWMutex）— 多个 goroutine 可能同时保存/加载检查点
//   - nil 输入不 panic
//   - Load 不存在的 key 返回 nil, false（不 panic）
//   - Delete 不存在的 key 静默成功（幂等）
type CheckpointStore struct {
	store map[string]*model.Checkpoint // key: TaskID
	mu    sync.RWMutex
}

// NewCheckpointStore 创建检查点存储。
func NewCheckpointStore() *CheckpointStore {
	return &CheckpointStore{
		store: make(map[string]*model.Checkpoint),
	}
}

// Save 保存任务的检查点（覆盖旧的）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. cp == nil → 静默返回（不 panic）
// 2. cp.TaskID == "" → 静默返回（无效数据不存储）
// 3. cp.Progress 不在 [0.0, 1.0] → clamp 到合法范围后存储
// 4. 使用 mu.Lock() 保护写操作
func (cs *CheckpointStore) Save(cp *model.Checkpoint) {
	// TODO: 实现
	panic("not implemented")
}

// Load 加载任务的最近检查点。
// 如果没有检查点，返回 nil, false。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. taskID == "" → 返回 nil, false
// 2. 使用 mu.RLock() 保护读操作
// 3. 返回检查点的副本（深拷贝），防止调用方修改内部状态
func (cs *CheckpointStore) Load(taskID string) (*model.Checkpoint, bool) {
	// TODO: 实现
	panic("not implemented")
}

// Delete 删除任务的检查点（任务完成或失败后清理）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. taskID == "" → 静默返回
// 2. key 不存在 → 静默返回（幂等操作）
// 3. 使用 mu.Lock() 保护写操作
func (cs *CheckpointStore) Delete(taskID string) {
	// TODO: 实现
	panic("not implemented")
}

// Count 返回当前存储的检查点数量。
func (cs *CheckpointStore) Count() int {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return len(cs.store)
}
