package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

// ============================================================
// 基础单元测试
// ============================================================

func TestAvailableMemory(t *testing.T) {
	// TODO(learner): 实现以下测试
	// 提示：创建 GPU_Node，设置 MemoryTotal 和 MemoryUsed，验证返回值

	t.Run("normal case", func(t *testing.T) {
		node := &GPU_Node{MemoryTotal: 8000, MemoryUsed: 3000}
		assert.Equal(t, 5000, node.AvailableMemory())
	})

	t.Run("fully used", func(t *testing.T) {
		// TODO: MemoryUsed == MemoryTotal → 返回 0
	})

	t.Run("overcommit returns zero not negative", func(t *testing.T) {
		// TODO: MemoryUsed > MemoryTotal（异常状态）→ 返回 0（防御性）
	})

	t.Run("zero total", func(t *testing.T) {
		// TODO: MemoryTotal == 0, MemoryUsed == 0 → 返回 0
	})
}

// ============================================================
// 表驱动测试
// ============================================================

func TestCanAccept(t *testing.T) {
	tests := []struct {
		name string
		node *GPU_Node
		req  ResourceRequirement
		want bool
	}{
		{
			name: "idle node with enough resources",
			node: &GPU_Node{Status: NodeStatusIdle, ComputePower: 100, MemoryTotal: 8000, MemoryUsed: 2000},
			req:  ResourceRequirement{ComputePower: 50, Memory: 4000},
			want: true,
		},
		{
			name: "offline node always rejects",
			node: &GPU_Node{Status: NodeStatusOffline, ComputePower: 100, MemoryTotal: 8000},
			req:  ResourceRequirement{ComputePower: 1, Memory: 1},
			want: false,
		},
		{
			name: "degraded node halves compute power - reject",
			node: &GPU_Node{Status: NodeStatusDegraded, ComputePower: 100, MemoryTotal: 8000},
			req:  ResourceRequirement{ComputePower: 60, Memory: 1000},
			want: false, // 100/2=50 < 60
		},
		{
			name: "degraded node halves compute power - accept",
			node: &GPU_Node{Status: NodeStatusDegraded, ComputePower: 100, MemoryTotal: 8000},
			req:  ResourceRequirement{ComputePower: 40, Memory: 1000},
			want: true, // 100/2=50 >= 40
		},
		// TODO(learner): 添加更多测试用例
		// - insufficient memory
		// - negative resource request (should return false)
		// - zero resource request
		// - busy node with remaining resources
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.node.CanAccept(tt.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ============================================================
// 错误路径测试
// ============================================================

func TestValidate(t *testing.T) {
	t.Run("valid node", func(t *testing.T) {
		node := &GPU_Node{
			ID:           "node-001",
			ComputePower: 100,
			MemoryTotal:  8000,
			FaultRate:    0.01,
		}
		assert.NoError(t, node.Validate())
	})

	t.Run("empty ID", func(t *testing.T) {
		node := &GPU_Node{ID: "", ComputePower: 100, MemoryTotal: 8000}
		err := node.Validate()
		assert.ErrorIs(t, err, ErrInvalidNodeID)
	})

	// TODO(learner): 添加更多验证测试
	// - negative ComputePower → ErrNegativeResource
	// - negative MemoryTotal → ErrNegativeResource
	// - MemoryUsed > MemoryTotal → ErrMemoryOvercommit
	// - FaultRate > 1.0 → ErrFaultRateOutOfRange
	// - FaultRate < 0.0 → ErrFaultRateOutOfRange
}

func TestAllocateTask(t *testing.T) {
	t.Run("successful allocation", func(t *testing.T) {
		// TODO(learner): 创建节点，分配任务，验证 MemoryUsed 和 AssignedTasks 更新
	})

	t.Run("idempotent - duplicate allocation", func(t *testing.T) {
		// TODO(learner): 同一任务分配两次，资源不应重复扣除
	})

	t.Run("insufficient resources", func(t *testing.T) {
		// TODO(learner): 资源不足时返回错误
	})

	t.Run("empty task ID", func(t *testing.T) {
		// TODO(learner): 空任务 ID 返回错误
	})
}

func TestReleaseTask(t *testing.T) {
	t.Run("successful release", func(t *testing.T) {
		// TODO(learner): 分配后释放，验证资源恢复
	})

	t.Run("release non-existent task", func(t *testing.T) {
		// TODO(learner): 释放不存在的任务返回错误
	})

	t.Run("memory does not go negative", func(t *testing.T) {
		// TODO(learner): 释放后 MemoryUsed 不应为负
	})
}

// ============================================================
// 属性测试（Property-Based Testing）
// ============================================================

func TestGPUNode_JSONRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机 GPU_Node
		node := &GPU_Node{
			ID:           rapid.StringMatching(`node-[0-9]{4}`).Draw(t, "id"),
			Status:       NodeStatus(rapid.IntRange(0, 3).Draw(t, "status")),
			ComputePower: rapid.IntRange(0, 10000).Draw(t, "power"),
			MemoryTotal:  rapid.IntRange(0, 1000000).Draw(t, "memTotal"),
			MemoryUsed:   rapid.IntRange(0, 1000000).Draw(t, "memUsed"),
			FaultRate:    rapid.Float64Range(0, 1).Draw(t, "faultRate"),
		}

		// 序列化
		data, err := json.Marshal(node)
		require.NoError(t, err, "marshal should not fail")

		// 反序列化
		var restored GPU_Node
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err, "unmarshal should not fail")

		// 验证往返一致性
		assert.Equal(t, node.ID, restored.ID)
		assert.Equal(t, node.Status, restored.Status)
		assert.Equal(t, node.ComputePower, restored.ComputePower)
		assert.Equal(t, node.MemoryTotal, restored.MemoryTotal)
		assert.Equal(t, node.MemoryUsed, restored.MemoryUsed)
		assert.InDelta(t, node.FaultRate, restored.FaultRate, 1e-10)
	})
}

func TestAvailableMemory_AlwaysNonNegative(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		node := &GPU_Node{
			MemoryTotal: rapid.IntRange(0, 100000).Draw(t, "total"),
			MemoryUsed:  rapid.IntRange(0, 200000).Draw(t, "used"), // 可能超过 total
		}

		result := node.AvailableMemory()

		// 属性：AvailableMemory 永远 >= 0
		assert.GreaterOrEqual(t, result, 0,
			"AvailableMemory must never be negative, got %d (total=%d, used=%d)",
			result, node.MemoryTotal, node.MemoryUsed)
	})
}
