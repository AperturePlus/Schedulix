package functools

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"schedulix/pkg/model"
)

// ============================================================
// Filter / Map / Reduce 测试
// ============================================================

func makeNodes(statuses ...model.NodeStatus) []*model.GPU_Node {
	nodes := make([]*model.GPU_Node, len(statuses))
	for i, s := range statuses {
		nodes[i] = &model.GPU_Node{
			ID: fmt.Sprintf("node-%d", i),
			Status: s,
			ComputePower: (i + 1) * 100,
			MemoryTotal: 8000,
			MemoryUsed: i * 1000,
		}
	}
	return nodes
}

func TestFilterNodes(t *testing.T) {
	nodes := makeNodes(model.NodeStatusIdle, model.NodeStatusOffline, model.NodeStatusIdle, model.NodeStatusBusy)

	idle := FilterNodes(nodes, func(n *model.GPU_Node) bool {
		return n.Status == model.NodeStatusIdle
	})

	assert.Len(t, idle, 2)
}

func TestFilterNodes_NilPredicate(t *testing.T) {
	nodes := makeNodes(model.NodeStatusIdle, model.NodeStatusOffline)
	result := FilterNodes(nodes, nil)
	assert.Len(t, result, 2) // 不过滤，返回全部
}

func TestFilterNodes_NilSlice(t *testing.T) {
	result := FilterNodes(nil, func(n *model.GPU_Node) bool { return true })
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestMapNodeScores(t *testing.T) {
	nodes := makeNodes(model.NodeStatusIdle, model.NodeStatusIdle)
	scores := MapNodeScores(nodes, func(n *model.GPU_Node) float64 {
		return float64(n.ComputePower)
	})

	assert.Equal(t, float64(100), scores["node-0"])
	assert.Equal(t, float64(200), scores["node-1"])
}

func TestMapNodeScores_ScorerPanic(t *testing.T) {
	nodes := makeNodes(model.NodeStatusIdle)
	scores := MapNodeScores(nodes, func(n *model.GPU_Node) float64 {
		panic("scorer panic!")
	})
	// panic 被 recover，得分为 0
	assert.Equal(t, float64(0), scores["node-0"])
}

func TestReduceNodes(t *testing.T) {
	nodes := makeNodes(model.NodeStatusIdle, model.NodeStatusIdle, model.NodeStatusIdle)
	total := ReduceNodes(nodes, 0, func(acc int, n *model.GPU_Node) int {
		return acc + n.MemoryTotal
	})
	assert.Equal(t, 24000, total) // 3 * 8000
}

// ============================================================
// 函数组合测试
// ============================================================

func TestComposePredicates(t *testing.T) {
	isIdle := func(n *model.GPU_Node) bool { return n.Status == model.NodeStatusIdle }
	hasMemory := func(n *model.GPU_Node) bool { return n.AvailableMemory() > 2000 }

	combined := ComposePredicates(isIdle, hasMemory)

	node := &model.GPU_Node{Status: model.NodeStatusIdle, MemoryTotal: 8000, MemoryUsed: 5000}
	assert.True(t, combined(node)) // Idle + 3000 > 2000

	node2 := &model.GPU_Node{Status: model.NodeStatusIdle, MemoryTotal: 8000, MemoryUsed: 7000}
	assert.False(t, combined(node2)) // Idle 但 1000 < 2000
}

func TestComposePredicates_Empty(t *testing.T) {
	combined := ComposePredicates()
	node := &model.GPU_Node{}
	assert.True(t, combined(node)) // 空组合 → 总是 true
}

func TestNegatePredicate(t *testing.T) {
	isOffline := func(n *model.GPU_Node) bool { return n.Status == model.NodeStatusOffline }
	notOffline := NegatePredicate(isOffline)

	assert.True(t, notOffline(&model.GPU_Node{Status: model.NodeStatusIdle}))
	assert.False(t, notOffline(&model.GPU_Node{Status: model.NodeStatusOffline}))
}

// ============================================================
// Pipeline 测试
// ============================================================

func TestNodePipeline(t *testing.T) {
	nodes := makeNodes(
		model.NodeStatusIdle, model.NodeStatusOffline,
		model.NodeStatusIdle, model.NodeStatusIdle,
		model.NodeStatusBusy,
	)

	result := NewNodePipeline().
		Filter(func(n *model.GPU_Node) bool { return n.Status == model.NodeStatusIdle }).
		SortBy(func(a, b *model.GPU_Node) bool { return a.ComputePower > b.ComputePower }).
		Limit(2).
		Execute(nodes)

	assert.Len(t, result, 2)
	// 排序后算力最高的两个 Idle 节点
	assert.True(t, result[0].ComputePower >= result[1].ComputePower)
}

func TestNodePipeline_NilInput(t *testing.T) {
	result := NewNodePipeline().Filter(func(n *model.GPU_Node) bool { return true }).Execute(nil)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

// ============================================================
// 函数选项模式测试
// ============================================================

func TestApplyOptions_Defaults(t *testing.T) {
	config := ApplyOptions()
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 4, config.WorkerCount)
	assert.True(t, config.EnableMetrics)
}

func TestApplyOptions_Custom(t *testing.T) {
	config := ApplyOptions(
		WithMaxRetries(10),
		WithWorkerCount(16),
		WithMetrics(false),
	)
	assert.Equal(t, 10, config.MaxRetries)
	assert.Equal(t, 16, config.WorkerCount)
	assert.False(t, config.EnableMetrics)
}

func TestApplyOptions_NilOption(t *testing.T) {
	config := ApplyOptions(nil, WithMaxRetries(5), nil)
	assert.Equal(t, 5, config.MaxRetries) // nil 被跳过
}

// ============================================================
// 闭包测试
// ============================================================

func TestMakeCounter(t *testing.T) {
	inc, get := MakeCounter()
	assert.Equal(t, 0, get())
	inc()
	inc()
	inc()
	assert.Equal(t, 3, get())
}

func TestMakeRetrier(t *testing.T) {
	callCount := 0
	retry := MakeRetrier(3, time.Millisecond)

	err := retry(func() error {
		callCount++
		if callCount < 3 {
			return errors.New("not yet")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, callCount) // 前两次失败，第三次成功
}

func TestMakeRetrier_AllFail(t *testing.T) {
	retry := MakeRetrier(3, time.Millisecond)
	err := retry(func() error {
		return errors.New("always fail")
	})
	assert.Error(t, err)
}

func TestMakeRetrier_NilFunc(t *testing.T) {
	retry := MakeRetrier(3, time.Millisecond)
	err := retry(nil)
	assert.Error(t, err)
}

// ============================================================
// HTTP 中间件测试
// ============================================================

func TestRecoveryMiddleware(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic!")
	})

	wrapped := RecoveryMiddleware()(panicHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req) // 不应 panic

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestChain(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建两个中间件，分别在 header 中添加 "X-A" 和 "X-B"
	// 2. Chain 它们
	// 3. 验证响应 header 包含两个值
	// 4. 验证执行顺序正确
}

func TestRateLimitMiddleware(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建限流中间件（每秒 2 次）
	// 2. 快速发送 3 个请求
	// 3. 前 2 个 → 200
	// 4. 第 3 个 → 429
}

// --- 防止 unused import ---
import "fmt"
var _ = fmt.Sprintf
var _ = require.NoError
