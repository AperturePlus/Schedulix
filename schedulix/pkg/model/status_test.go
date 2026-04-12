package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestNodeStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		status NodeStatus
		want   string
	}{
		{NodeStatusIdle, `"idle"`},
		{NodeStatusBusy, `"busy"`},
		{NodeStatusOffline, `"offline"`},
		{NodeStatusDegraded, `"degraded"`},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			data, err := json.Marshal(tt.status)
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(data))
		})
	}
}

func TestNodeStatus_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input string
		want  NodeStatus
	}{
		{`"idle"`, NodeStatusIdle},
		{`"busy"`, NodeStatusBusy},
		{`"offline"`, NodeStatusOffline},
		{`"degraded"`, NodeStatusDegraded},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var got NodeStatus
			err := json.Unmarshal([]byte(tt.input), &got)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNodeStatus_UnmarshalJSON_InvalidInput(t *testing.T) {
	// TODO(learner): 测试不合法的输入
	// - 未知状态字符串 → 返回错误
	// - 非字符串 JSON（如数字）→ 返回错误
	// - 空字符串 → 返回错误

	t.Run("unknown status", func(t *testing.T) {
		var s NodeStatus
		err := json.Unmarshal([]byte(`"unknown"`), &s)
		assert.Error(t, err)
	})

	t.Run("number instead of string", func(t *testing.T) {
		// TODO: 实现
	})
}

func TestNodeStatus_RoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		status := NodeStatus(rapid.IntRange(0, 3).Draw(t, "status"))

		data, err := json.Marshal(status)
		require.NoError(t, err)

		var restored NodeStatus
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		assert.Equal(t, status, restored)
	})
}

// TODO(learner): 为 TaskStatus 编写相同的测试
// - TestTaskStatus_MarshalJSON
// - TestTaskStatus_UnmarshalJSON
// - TestTaskStatus_UnmarshalJSON_InvalidInput
// - TestTaskStatus_RoundTrip
