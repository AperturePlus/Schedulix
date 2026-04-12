package scheduler

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"schedulix/pkg/model"
	"schedulix/pkg/queue"
)

// ============================================================
// 调度策略测试
// ============================================================

func makeCluster(n int, memTotal int) *model.Cluster {
	c := model.NewCluster(n)
	for _, node := range c.Nodes {
		node.MemoryTotal = memTotal
		node.ComputePower = 100
	}
	return c
}

func makeTask(id string, mem int) *model.Task {
	return &model.Task{
		ID:       id,
		Resource: model.ResourceRequirement{ComputePower: 10, Memory: mem},
		Priority: 5,
	}
}

func TestFirstFit_BasicSchedule(t *testing.T) {
	c := makeCluster(5, 8000)
	s := &FirstFitStrategy{}
	task := makeTask("task-1", 4000)

	nodeID, err := s.Schedule(task, c)
	require.NoError(t, err)
	assert.NotEmpty(t, nodeID)
}

func TestFirstFit_NoAvailableNode(t *testing.T) {
	c := makeCluster(3, 1000)
	s := &FirstFitStrategy{}
	task := makeTask("task-1", 5000) // 需要 5000，所有节点只有 1000

	_, err := s.Schedule(task, c)
	assert.ErrorIs(t, err, ErrNoAvailableNode)
}

func TestBestFit_SelectsSmallestFit(t *testing.T) {
	// TODO(learner): 实现
	// 创建 3 个节点，内存分别为 2000, 5000, 10000
	// 任务需要 1500
	// Best-Fit 应选择 2000 的节点（剩余最少）
}

func TestRoundRobin_DistributesEvenly(t *testing.T) {
	// TODO(learner): 实现
	// 创建 3 个节点
	// 调度 6 个任务
	// 每个节点应分到 2 个任务
}

// ============================================================
// 输入验证测试
// ============================================================

func TestValidateScheduleInput(t *testing.T) {
	c := makeCluster(5, 8000)

	t.Run("nil task", func(t *testing.T) {
		err := ValidateScheduleInput(nil, c)
		assert.ErrorIs(t, err, ErrNilTask)
	})

	t.Run("nil cluster", func(t *testing.T) {
		task := makeTask("t-1", 1000)
		err := ValidateScheduleInput(task, nil)
		assert.ErrorIs(t, err, ErrNilCluster)
	})

	t.Run("negative resource", func(t *testing.T) {
		task := &model.Task{
			ID:       "t-1",
			Resource: model.ResourceRequirement{ComputePower: -1, Memory: 1000},
		}
		err := ValidateScheduleInput(task, c)
		assert.ErrorIs(t, err, ErrInvalidResourceRequest)
	})

	t.Run("already assigned", func(t *testing.T) {
		task := makeTask("t-1", 1000)
		task.AssignedNodeID = "node-0001"
		err := ValidateScheduleInput(task, c)
		assert.ErrorIs(t, err, ErrTaskAlreadyAssigned)
	})

	t.Run("valid input", func(t *testing.T) {
		task := makeTask("t-1", 1000)
		err := ValidateScheduleInput(task, c)
		assert.NoError(t, err)
	})
}

// ============================================================
// Scheduler 集成测试
// ============================================================

func TestScheduler_ScheduleNext(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建集群、队列、调度器
	// 2. 入队一个任务
	// 3. 调用 ScheduleNext
	// 4. 验证任务状态变为 Running
	// 5. 验证节点资源已更新
}

func TestScheduler_ScheduleNext_EmptyQueue(t *testing.T) {
	c := makeCluster(5, 8000)
	q := queue.NewTaskQueue()
	s := NewScheduler(&FirstFitStrategy{}, q, c)

	task, err := s.ScheduleNext()
	assert.Nil(t, task)
	assert.ErrorIs(t, err, queue.ErrQueueEmpty)
}

func TestScheduler_TaskNotLost_OnFailure(t *testing.T) {
	// TODO(learner): 关键测试！
	// 1. 创建一个资源很少的集群（所有节点内存 = 100）
	// 2. 入队一个需要大量资源的任务（内存 = 10000）
	// 3. 调用 ScheduleNext → 应该失败
	// 4. 验证任务仍在队列中（没有丢失！）
}

// ============================================================
// 并发调度测试
// ============================================================

func TestConcurrentScheduler_NoResourceOvercommit(t *testing.T) {
	// TODO(learner): 关键测试！
	// 1. 创建 10 个节点，每个 1000MB 内存
	// 2. 提交 100 个任务，每个需要 500MB
	// 3. 最多只能调度 20 个任务（10 节点 × 2 任务/节点）
	// 4. 验证没有节点的 MemoryUsed > MemoryTotal
	// 5. 必须用 go test -race 运行
}

func TestConcurrentScheduler_SubmitTimeout(t *testing.T) {
	cs := NewConcurrentScheduler(&FirstFitStrategy{}, makeCluster(5, 8000))
	defer cs.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	err := cs.Submit(ctx, makeTask("t-1", 1000))
	assert.Error(t, err) // 应该返回超时/取消错误
}

func TestConcurrentScheduler_StopIdempotent(t *testing.T) {
	cs := NewConcurrentScheduler(&FirstFitStrategy{}, makeCluster(5, 8000))

	// 多次 Stop 不应 panic
	cs.Stop()
	cs.Stop()
	cs.Stop()
}

func TestConcurrentScheduler_SubmitAfterStop(t *testing.T) {
	cs := NewConcurrentScheduler(&FirstFitStrategy{}, makeCluster(5, 8000))
	cs.Stop()

	err := cs.Submit(context.Background(), makeTask("t-1", 1000))
	assert.ErrorIs(t, err, ErrChannelClosed)
}

// ============================================================
// 基准测试
// ============================================================

func BenchmarkFirstFit_100Nodes(b *testing.B) {
	c := makeCluster(100, 80000)
	s := &FirstFitStrategy{}
	task := makeTask("bench", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Schedule(task, c)
	}
}

func BenchmarkFirstFit_10000Nodes(b *testing.B) {
	c := makeCluster(10000, 80000)
	s := &FirstFitStrategy{}
	task := makeTask("bench", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Schedule(task, c)
	}
}

func BenchmarkBestFit_10000Nodes(b *testing.B) {
	c := makeCluster(10000, 80000)
	s := &BestFitStrategy{}
	task := makeTask("bench", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Schedule(task, c)
	}
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
var _ sync.WaitGroup
