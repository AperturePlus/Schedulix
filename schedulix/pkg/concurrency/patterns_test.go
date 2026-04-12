package concurrency

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Fan-Out / Fan-In 测试
// ============================================================

func TestFanOut(t *testing.T) {
	ctx := context.Background()
	inputs := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		inputs <- i
	}
	close(inputs)

	// 3 个 worker 并行计算平方
	results := FanOut(ctx, inputs, 3, func(n int) int {
		return n * n
	})

	var sum int
	for r := range results {
		sum += r
	}
	// 1² + 2² + ... + 10² = 385
	assert.Equal(t, 385, sum)
}

func TestFanOut_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	inputs := make(chan int) // 无缓冲，会阻塞

	results := FanOut(ctx, inputs, 3, func(n int) int {
		return n
	})

	cancel() // 取消

	// results channel 应该最终关闭
	count := 0
	for range results {
		count++
	}
	// 可能收到 0 个或少量结果
	t.Logf("received %d results after cancel", count)
}

func TestFanOut_WorkerPanic(t *testing.T) {
	ctx := context.Background()
	inputs := make(chan int, 3)
	inputs <- 1
	inputs <- 2
	inputs <- 3
	close(inputs)

	callCount := atomic.Int32{}
	results := FanOut(ctx, inputs, 2, func(n int) int {
		callCount.Add(1)
		if n == 2 {
			panic("boom!")
		}
		return n * 10
	})

	var got []int
	for r := range results {
		got = append(got, r)
	}
	// n=2 的 panic 被 recover，其他结果正常
	assert.Contains(t, got, 10)  // n=1 → 10
	assert.Contains(t, got, 30)  // n=3 → 30
}

func TestFanIn(t *testing.T) {
	ctx := context.Background()

	ch1 := make(chan int, 2)
	ch2 := make(chan int, 2)
	ch1 <- 1
	ch1 <- 2
	close(ch1)
	ch2 <- 3
	ch2 <- 4
	close(ch2)

	merged := FanIn(ctx, ch1, ch2)

	var got []int
	for v := range merged {
		got = append(got, v)
	}
	assert.Len(t, got, 4)
	assert.ElementsMatch(t, []int{1, 2, 3, 4}, got)
}

// ============================================================
// Worker Pool 测试
// ============================================================

func TestWorkerPool_BasicUsage(t *testing.T) {
	pool := NewWorkerPool[int, int](3, 10, func(n int) (int, error) {
		return n * 2, nil
	})

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		pool.Submit(ctx, Job[int, int]{ID: fmt.Sprintf("j-%d", i), Input: i})
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		pool.Shutdown()
	}()

	var results []int
	for r := range pool.Results() {
		require.NoError(t, r.Err)
		results = append(results, r.Value)
	}
	assert.Len(t, results, 5)
}

func TestWorkerPool_SubmitAfterShutdown(t *testing.T) {
	pool := NewWorkerPool[int, int](2, 5, func(n int) (int, error) {
		return n, nil
	})
	pool.Shutdown()

	err := pool.Submit(context.Background(), Job[int, int]{ID: "late", Input: 1})
	assert.ErrorIs(t, err, ErrWorkerPoolClosed)
}

func TestWorkerPool_WorkerPanic(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建 pool，process 函数在某些输入时 panic
	// 2. 提交多个任务
	// 3. 验证 panic 的任务返回 error，其他任务正常
}

// ============================================================
// Pipeline 测试
// ============================================================

func TestPipelineStage(t *testing.T) {
	ctx := context.Background()

	// 阶段 1：生成数字
	gen := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		gen <- i
	}
	close(gen)

	// 阶段 2：乘以 2
	double := PipelineStage(func(n int) int { return n * 2 })
	doubled := double(ctx, gen)

	// 阶段 3：转为字符串
	toString := PipelineStage(func(n int) string { return fmt.Sprintf("val-%d", n) })
	strings := toString(ctx, doubled)

	var got []string
	for s := range strings {
		got = append(got, s)
	}
	assert.ElementsMatch(t, []string{"val-2", "val-4", "val-6", "val-8", "val-10"}, got)
}

// ============================================================
// Select / Timeout 测试
// ============================================================

func TestFirstOf(t *testing.T) {
	ctx := context.Background()

	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	// ch2 先产生结果
	ch2 <- "fast"

	result, err := FirstOf(ctx, ch1, ch2)
	require.NoError(t, err)
	assert.Equal(t, "fast", result)
}

func TestFirstOf_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	ch1 := make(chan string) // 永远不产生结果

	_, err := FirstOf(ctx, ch1)
	assert.Error(t, err) // 应该超时
}

func TestTimeout_Success(t *testing.T) {
	result, err := Timeout(time.Second, func() (int, error) {
		return 42, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestTimeout_Exceeded(t *testing.T) {
	_, err := Timeout(50*time.Millisecond, func() (int, error) {
		time.Sleep(time.Second)
		return 0, nil
	})
	assert.ErrorIs(t, err, ErrTimeout)
}

func TestTimeout_FuncPanic(t *testing.T) {
	_, err := Timeout(time.Second, func() (int, error) {
		panic("boom!")
	})
	assert.Error(t, err)
}

// ============================================================
// LazyInit 测试
// ============================================================

func TestLazyInit(t *testing.T) {
	callCount := 0
	lazy := NewLazyInit(func() (string, error) {
		callCount++
		return "hello", nil
	})

	v1, err := lazy.Get()
	require.NoError(t, err)
	assert.Equal(t, "hello", v1)

	v2, _ := lazy.Get()
	assert.Equal(t, "hello", v2)
	assert.Equal(t, 1, callCount) // 只调用了一次
}

func TestLazyInit_Concurrent(t *testing.T) {
	callCount := atomic.Int32{}
	lazy := NewLazyInit(func() (int, error) {
		callCount.Add(1)
		time.Sleep(10 * time.Millisecond)
		return 42, nil
	})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, _ := lazy.Get()
			assert.Equal(t, 42, v)
		}()
	}
	wg.Wait()
	assert.Equal(t, int32(1), callCount.Load()) // 100 个 goroutine，只初始化一次
}

// ============================================================
// Semaphore 测试
// ============================================================

func TestSemaphore(t *testing.T) {
	sem := NewSemaphore(3)
	var running atomic.Int32
	var maxRunning atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem.Acquire(context.Background())
			defer sem.Release()

			cur := running.Add(1)
			// 记录最大并发数
			for {
				old := maxRunning.Load()
				if cur <= old || maxRunning.CompareAndSwap(old, cur) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
			running.Add(-1)
		}()
	}
	wg.Wait()

	assert.LessOrEqual(t, maxRunning.Load(), int32(3)) // 最多 3 个并发
}

func TestSemaphore_ContextCancel(t *testing.T) {
	sem := NewSemaphore(1)
	sem.Acquire(context.Background()) // 占满

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := sem.Acquire(ctx) // 应该超时
	assert.Error(t, err)

	sem.Release() // 释放
}

// ============================================================
// ErrGroup 测试
// ============================================================

func TestErrGroup_AllSuccess(t *testing.T) {
	g := NewErrGroup(context.Background())
	g.Go(func(ctx context.Context) error { return nil })
	g.Go(func(ctx context.Context) error { return nil })
	assert.NoError(t, g.Wait())
}

func TestErrGroup_FirstError(t *testing.T) {
	g := NewErrGroup(context.Background())
	g.Go(func(ctx context.Context) error {
		return errors.New("fail")
	})
	g.Go(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})
	assert.Error(t, g.Wait())
}

func TestErrGroup_CancelsOnError(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建 ErrGroup
	// 2. goroutine A 立即返回错误
	// 3. goroutine B 检查 ctx.Done() → 应该被取消
	// 4. 验证 B 收到了取消信号
}

func TestErrGroup_PanicRecovery(t *testing.T) {
	// TODO(learner): 实现
	// 1. goroutine panic
	// 2. Wait() 应返回 error（不 panic）
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
