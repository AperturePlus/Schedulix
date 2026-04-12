package queue

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"schedulix/pkg/model"
)

// ============================================================
// 基础功能测试
// ============================================================

func TestEnqueueDequeue_Priority(t *testing.T) {
	q := NewTaskQueue()

	// 入队三个不同优先级的任务
	q.Enqueue(&model.Task{ID: "low", Priority: 1})
	q.Enqueue(&model.Task{ID: "high", Priority: 10})
	q.Enqueue(&model.Task{ID: "mid", Priority: 5})

	// 出队应按优先级从高到低
	task, err := q.Dequeue()
	require.NoError(t, err)
	assert.Equal(t, "high", task.ID)

	task, err = q.Dequeue()
	require.NoError(t, err)
	assert.Equal(t, "mid", task.ID)

	task, err = q.Dequeue()
	require.NoError(t, err)
	assert.Equal(t, "low", task.ID)
}

func TestEnqueueDequeue_SamePriority_FIFO(t *testing.T) {
	q := NewTaskQueue()
	now := time.Now()

	// 同优先级，按提交时间 FIFO
	q.Enqueue(&model.Task{ID: "first", Priority: 5, SubmitTime: now})
	q.Enqueue(&model.Task{ID: "second", Priority: 5, SubmitTime: now.Add(time.Second)})
	q.Enqueue(&model.Task{ID: "third", Priority: 5, SubmitTime: now.Add(2 * time.Second)})

	task, _ := q.Dequeue()
	assert.Equal(t, "first", task.ID)
	task, _ = q.Dequeue()
	assert.Equal(t, "second", task.ID)
	task, _ = q.Dequeue()
	assert.Equal(t, "third", task.ID)
}

// ============================================================
// 错误路径测试
// ============================================================

func TestDequeue_EmptyQueue(t *testing.T) {
	q := NewTaskQueue()
	task, err := q.Dequeue()
	assert.Nil(t, task)
	assert.True(t, errors.Is(err, ErrQueueEmpty))
}

func TestPeek_EmptyQueue(t *testing.T) {
	q := NewTaskQueue()
	task, err := q.Peek()
	assert.Nil(t, task)
	assert.True(t, errors.Is(err, ErrQueueEmpty))
}

func TestEnqueue_NilTask(t *testing.T) {
	q := NewTaskQueue()
	err := q.Enqueue(nil)
	assert.True(t, errors.Is(err, ErrNilTask))
}

func TestEnqueue_EmptyTaskID(t *testing.T) {
	q := NewTaskQueue()
	err := q.Enqueue(&model.Task{ID: ""})
	assert.True(t, errors.Is(err, ErrInvalidTask))
}

func TestEnqueue_QueueFull(t *testing.T) {
	q := NewTaskQueueWithCapacity(3)
	q.Enqueue(&model.Task{ID: "1", Priority: 1})
	q.Enqueue(&model.Task{ID: "2", Priority: 2})
	q.Enqueue(&model.Task{ID: "3", Priority: 3})

	err := q.Enqueue(&model.Task{ID: "4", Priority: 4})
	assert.True(t, errors.Is(err, ErrQueueFull))
}

// ============================================================
// Len / IsEmpty 测试
// ============================================================

func TestLenAndIsEmpty(t *testing.T) {
	q := NewTaskQueue()
	assert.True(t, q.IsEmpty())
	assert.Equal(t, 0, q.Len())

	q.Enqueue(&model.Task{ID: "1", Priority: 1})
	assert.False(t, q.IsEmpty())
	assert.Equal(t, 1, q.Len())

	q.Dequeue()
	assert.True(t, q.IsEmpty())
	assert.Equal(t, 0, q.Len())
}

// ============================================================
// Drain 测试
// ============================================================

func TestDrain(t *testing.T) {
	// TODO(learner): 实现
	// 1. 入队多个任务
	// 2. 调用 Drain()
	// 3. 验证返回的任务按优先级排序
	// 4. 验证队列为空
}

// ============================================================
// 并发测试（必须用 go test -race 运行）
// ============================================================

func TestConcurrentEnqueueDequeue(t *testing.T) {
	q := NewTaskQueue()
	const n = 1000
	var wg sync.WaitGroup

	// 并发入队
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			q.Enqueue(&model.Task{
				ID:       fmt.Sprintf("task-%d", id),
				Priority: id % 10,
			})
		}(i)
	}
	wg.Wait()
	assert.Equal(t, n, q.Len())

	// 并发出队
	results := make(chan *model.Task, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task, err := q.Dequeue()
			if err == nil {
				results <- task
			}
		}()
	}
	wg.Wait()
	close(results)

	// 验证无重复
	seen := make(map[string]bool)
	for task := range results {
		assert.False(t, seen[task.ID], "duplicate task: %s", task.ID)
		seen[task.ID] = true
	}
}

// ============================================================
// 基准测试
// ============================================================

func BenchmarkEnqueue(b *testing.B) {
	q := NewTaskQueue()
	task := &model.Task{ID: "bench", Priority: 5}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(&model.Task{ID: fmt.Sprintf("t-%d", i), Priority: task.Priority})
	}
}

func BenchmarkDequeue(b *testing.B) {
	q := NewTaskQueue()
	for i := 0; i < b.N; i++ {
		q.Enqueue(&model.Task{ID: fmt.Sprintf("t-%d", i), Priority: i % 100})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Dequeue()
	}
}
