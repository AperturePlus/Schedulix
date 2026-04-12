package simulator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"

	"schedulix/pkg/model"
)

// ============================================================
// EventConfig 测试
// ============================================================

func TestEventConfig_Validate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		c := &EventConfig{
			NodeDownProb: 0.01, NetworkDelayProb: 0.02,
			DegradedProb: 0.01, RecoveryProb: 0.05,
			TotalSteps: 100, StepIntervalMs: 100,
		}
		assert.NoError(t, c.Validate())
	})

	t.Run("probability out of range", func(t *testing.T) {
		c := &EventConfig{NodeDownProb: 1.5, TotalSteps: 100, StepIntervalMs: 100}
		assert.Error(t, c.Validate())
	})

	t.Run("negative probability", func(t *testing.T) {
		c := &EventConfig{NodeDownProb: -0.1, TotalSteps: 100, StepIntervalMs: 100}
		assert.Error(t, c.Validate())
	})

	t.Run("zero total steps", func(t *testing.T) {
		c := DefaultEventConfig()
		c.TotalSteps = 0
		assert.Error(t, c.Validate())
	})

	// TODO(learner): 添加更多测试
	// - 多个字段同时不合法 → 错误消息应包含所有问题
	// - StepIntervalMs <= 0
}

func TestParseConfig(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		data := []byte(`{"node_down_prob":0.01,"network_delay_prob":0.02,"degraded_prob":0.01,"recovery_prob":0.05,"total_steps":100,"step_interval_ms":100}`)
		cfg, err := ParseConfig(data)
		require.NoError(t, err)
		assert.Equal(t, 0.01, cfg.NodeDownProb)
		assert.Equal(t, 100, cfg.TotalSteps)
	})

	t.Run("empty input returns default", func(t *testing.T) {
		cfg, err := ParseConfig(nil)
		assert.Error(t, err) // 有警告
		assert.NotNil(t, cfg) // 但返回默认配置
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := ParseConfig([]byte("not json"))
		assert.Error(t, err)
	})
}

func TestEventConfig_JSONRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		cfg := &EventConfig{
			NodeDownProb:     rapid.Float64Range(0, 1).Draw(t, "ndp"),
			NetworkDelayProb: rapid.Float64Range(0, 1).Draw(t, "nlp"),
			DegradedProb:     rapid.Float64Range(0, 1).Draw(t, "dp"),
			RecoveryProb:     rapid.Float64Range(0, 1).Draw(t, "rp"),
			TotalSteps:       rapid.IntRange(1, 10000).Draw(t, "steps"),
			StepIntervalMs:   int64(rapid.IntRange(1, 10000).Draw(t, "interval")),
		}

		data, err := json.Marshal(cfg)
		require.NoError(t, err)

		var restored EventConfig
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		assert.InDelta(t, cfg.NodeDownProb, restored.NodeDownProb, 1e-10)
		assert.Equal(t, cfg.TotalSteps, restored.TotalSteps)
	})
}

// ============================================================
// EventSimulator 测试
// ============================================================

func TestEventSimulator_NilConfig(t *testing.T) {
	cluster := model.NewCluster(10)
	sim := NewEventSimulator(nil, cluster)
	assert.NotNil(t, sim) // 应使用默认配置
}

func TestEventSimulator_NilCluster(t *testing.T) {
	cfg := DefaultEventConfig()
	sim := NewEventSimulator(cfg, nil)
	err := sim.RunStepMode()
	assert.Error(t, err) // 应返回 ErrSimulatorNotReady
}

func TestEventSimulator_RegisterNilHandler(t *testing.T) {
	sim := NewEventSimulator(DefaultEventConfig(), model.NewCluster(10))
	sim.RegisterHandler(nil) // 不应 panic
}

func TestEventSimulator_RunStepMode(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建集群（10 个节点）
	// 2. 配置高故障概率（如 0.5）以确保产生事件
	// 3. 运行 10 步
	// 4. 验证事件日志不为空
	// 5. 验证事件日志按时间排序
}

// ============================================================
// Handler Panic 隔离测试
// ============================================================

type panicHandler struct{}

func (h *panicHandler) OnFault(event *FaultEvent) error {
	panic("handler panic!")
}

func (h *panicHandler) OnRecovery(event *FaultEvent) error {
	return nil
}

func TestEventSimulator_HandlerPanicRecovery(t *testing.T) {
	// TODO(learner): 实现
	// 1. 注册一个会 panic 的 handler
	// 2. 运行模拟
	// 3. 验证模拟没有崩溃
	// 4. 验证错误被记录到 GetErrors()
}
