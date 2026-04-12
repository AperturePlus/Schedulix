package gateway

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"schedulix/pkg/model"
	"schedulix/pkg/queue"
	"schedulix/pkg/scheduler"
)

// ============================================================
// HTTP Handler 测试
// 使用 httptest 包模拟 HTTP 请求，不需要启动真实服务器
// ============================================================

func setupHandler(t *testing.T) *Handler {
	t.Helper()
	c := model.NewCluster(10)
	for _, node := range c.Nodes {
		node.MemoryTotal = 8000
		node.ComputePower = 100
	}
	q := queue.NewTaskQueue()
	s := scheduler.NewScheduler(&scheduler.FirstFitStrategy{}, q, c)
	return NewHandler(c, s, q)
}

func TestSubmitTask_ValidRequest(t *testing.T) {
	h := setupHandler(t)

	task := model.Task{
		ID:       "task-1",
		Priority: 5,
		Resource: model.ResourceRequirement{ComputePower: 10, Memory: 1000},
	}
	body, _ := json.Marshal(task)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.SubmitTask(w, req)

	// TODO(learner): 验证响应
	// assert.Equal(t, http.StatusAccepted, w.Code)
	// 解析响应 body，验证包含任务 ID
}

func TestSubmitTask_InvalidJSON(t *testing.T) {
	h := setupHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.SubmitTask(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSubmitTask_MissingFields(t *testing.T) {
	// TODO(learner): 实现
	// 提交缺少 ID 的任务 → 400
	// 提交 Memory <= 0 的任务 → 400
}

func TestGetClusterStatus_NilCluster(t *testing.T) {
	h := NewHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cluster/status", nil)
	w := httptest.NewRecorder()

	h.GetClusterStatus(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

// TODO(learner): 添加更多 handler 测试
// - TestGetNodes
// - TestGetMetrics
// - TestExportMetrics_JSON
// - TestExportMetrics_CSV
// - TestExportMetrics_InvalidFormat → 400

// ============================================================
// Scaler 测试
// ============================================================

func TestScaler_ColdStart(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建 Scaler（冷启动延迟 10ms）
	// 2. 第一次 OnRequest → isColdStart == true
	// 3. 第二次 OnRequest → isColdStart == false
}

func TestScaler_ScaleToZero(t *testing.T) {
	// TODO(learner): 实现
	// 1. OnRequest → activeInstances == 1
	// 2. OnRequestDone → activeInstances == 0
	// 3. 等待缩容延迟后，验证实例数为 0
}

func TestScaler_NegativeInstances(t *testing.T) {
	// TODO(learner): 实现
	// 多次 OnRequestDone（比 OnRequest 多）→ activeInstances 不应为负
}

// --- 防止 unused import ---
var _ = require.NoError
