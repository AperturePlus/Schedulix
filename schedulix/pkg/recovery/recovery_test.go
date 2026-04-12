package recovery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"schedulix/pkg/model"
	"schedulix/pkg/queue"
	"schedulix/pkg/simulator"
)

func TestCheckpointStore_SaveAndLoad(t *testing.T) {
	cs := NewCheckpointStore()

	cp := &model.Checkpoint{
		TaskID:   "task-1",
		Progress: 0.5,
		NodeID:   "node-1",
	}
	cs.Save(cp)

	loaded, ok := cs.Load("task-1")
	require.True(t, ok)
	assert.Equal(t, 0.5, loaded.Progress)
}

func TestCheckpointStore_LoadNonExistent(t *testing.T) {
	cs := NewCheckpointStore()
	_, ok := cs.Load("missing")
	assert.False(t, ok)
}

func TestCheckpointStore_SaveNil(t *testing.T) {
	cs := NewCheckpointStore()
	cs.Save(nil) // 不应 panic
	assert.Equal(t, 0, cs.Count())
}

func TestCheckpointStore_DeleteIdempotent(t *testing.T) {
	cs := NewCheckpointStore()
	cs.Delete("nonexistent") // 不应 panic
}

func TestCheckpointStore_ReturnsCopy(t *testing.T) {
	// TODO(learner): 实现
	// 1. Save 一个检查点
	// 2. Load 获取副本
	// 3. 修改副本的 Progress
	// 4. 再次 Load → Progress 应该不变（返回的是副本）
}

func TestRecoveryEngine_OnFault(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建集群，将一些任务分配到 node-0001
	// 2. 创建 RecoveryEngine
	// 3. 触发 node-0001 宕机事件
	// 4. 验证受影响任务被重新入队
	// 5. 验证 RecoveryLog 记录正确
}

func TestRecoveryEngine_MaxMigrations(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建一个 MigrationCount = 2 的任务
	// 2. 触发故障 → MigrationCount 变为 3
	// 3. 验证任务被标记为 Failed（不再入队）
}

func TestRecoveryEngine_IdempotentFault(t *testing.T) {
	// TODO(learner): 实现
	// 同一个故障事件处理两次，不应重复迁移任务
}

func TestRecoveryEngine_NilEvent(t *testing.T) {
	c := model.NewCluster(10)
	q := queue.NewTaskQueue()
	re := NewRecoveryEngine(c, q)

	err := re.OnFault(nil)
	assert.Error(t, err) // 不应 panic
}

// --- 防止 unused import ---
var _ = time.Now
var _ = simulator.FaultNodeDown
