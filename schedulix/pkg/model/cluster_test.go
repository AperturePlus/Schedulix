package model

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestNewCluster(t *testing.T) {
	t.Run("creates correct number of nodes", func(t *testing.T) {
		c := NewCluster(100)
		require.NotNil(t, c)
		assert.Len(t, c.Nodes, 100)
	})

	t.Run("all nodes have unique IDs", func(t *testing.T) {
		c := NewCluster(1000)
		ids := make(map[string]bool)
		for id := range c.Nodes {
			assert.False(t, ids[id], "duplicate ID: %s", id)
			ids[id] = true
		}
	})

	t.Run("all nodes start as Idle", func(t *testing.T) {
		c := NewCluster(50)
		for _, node := range c.Nodes {
			assert.Equal(t, NodeStatusIdle, node.Status)
		}
	})

	t.Run("zero count", func(t *testing.T) {
		// TODO(learner): count <= 0 应该如何处理？
		// 你的设计决策：返回 nil？返回空集群？panic？
	})

	t.Run("negative count", func(t *testing.T) {
		// TODO(learner): 同上
	})
}

func TestGetNode(t *testing.T) {
	c := NewCluster(10)

	t.Run("existing node", func(t *testing.T) {
		node, ok := c.GetNode("node-0001")
		assert.True(t, ok)
		assert.NotNil(t, node)
		assert.Equal(t, "node-0001", node.ID)
	})

	t.Run("non-existing node", func(t *testing.T) {
		node, ok := c.GetNode("node-9999")
		assert.False(t, ok)
		assert.Nil(t, node)
	})

	t.Run("empty ID", func(t *testing.T) {
		// TODO(learner): 空 ID 应返回 nil, false
	})
}

func TestUpdateNodeStatus(t *testing.T) {
	t.Run("update and verify", func(t *testing.T) {
		c := NewCluster(10)
		err := c.UpdateNodeStatus("node-0001", NodeStatusOffline)
		require.NoError(t, err)

		node, _ := c.GetNode("node-0001")
		assert.Equal(t, NodeStatusOffline, node.Status)
	})

	t.Run("status index consistency", func(t *testing.T) {
		// TODO(learner): 更新状态后，GetAvailableNodes 应反映变化
		// 1. 创建集群
		// 2. 将 node-0001 设为 Offline
		// 3. GetAvailableNodes(NodeStatusIdle) 不应包含 node-0001
		// 4. GetAvailableNodes(NodeStatusOffline) 应包含 node-0001
	})

	t.Run("idempotent - same status", func(t *testing.T) {
		// TODO(learner): 设置相同状态应该成功且无副作用
	})

	t.Run("non-existing node", func(t *testing.T) {
		c := NewCluster(10)
		err := c.UpdateNodeStatus("node-9999", NodeStatusOffline)
		assert.ErrorIs(t, err, ErrNodeNotFound)
	})
}

// ============================================================
// 并发测试
// ============================================================

func TestCluster_ConcurrentAccess(t *testing.T) {
	c := NewCluster(100)
	var wg sync.WaitGroup

	// 并发读写
	for i := 0; i < 100; i++ {
		wg.Add(2)

		// 并发读
		go func(id int) {
			defer wg.Done()
			nodeID := fmt.Sprintf("node-%04d", id+1)
			c.GetNode(nodeID)
			c.GetAvailableNodes(NodeStatusIdle)
		}(i)

		// 并发写
		go func(id int) {
			defer wg.Done()
			nodeID := fmt.Sprintf("node-%04d", id+1)
			status := NodeStatus(id % 4)
			c.UpdateNodeStatus(nodeID, status)
		}(i)
	}

	wg.Wait()
	// 如果没有 panic 或 data race，测试通过
	// 必须用 go test -race 运行！
}

// ============================================================
// 属性测试
// ============================================================

func TestCluster_SnapshotRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		count := rapid.IntRange(1, 100).Draw(t, "count")
		c := NewCluster(count)

		// 随机修改一些节点状态
		numChanges := rapid.IntRange(0, count).Draw(t, "changes")
		for i := 0; i < numChanges; i++ {
			nodeID := fmt.Sprintf("node-%04d", rapid.IntRange(1, count).Draw(t, "nodeIdx"))
			status := NodeStatus(rapid.IntRange(0, 3).Draw(t, "status"))
			c.UpdateNodeStatus(nodeID, status)
		}

		// 快照
		data, err := c.SnapshotToJSON()
		require.NoError(t, err)

		// 恢复
		restored, err := RestoreFromJSON(data)
		require.NoError(t, err)

		// 验证节点数量一致
		assert.Equal(t, len(c.Nodes), len(restored.Nodes))

		// 验证每个节点状态一致
		for id, node := range c.Nodes {
			rNode, ok := restored.Nodes[id]
			require.True(t, ok, "node %s missing after restore", id)
			assert.Equal(t, node.Status, rNode.Status)
			assert.Equal(t, node.ComputePower, rNode.ComputePower)
		}
	})
}

// ============================================================
// 基准测试
// ============================================================

func BenchmarkNewCluster_10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewCluster(10000)
	}
}

func BenchmarkGetAvailableNodes_10000(b *testing.B) {
	c := NewCluster(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetAvailableNodes(NodeStatusIdle)
	}
}

func BenchmarkSnapshotToJSON_10000(b *testing.B) {
	c := NewCluster(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.SnapshotToJSON()
	}
}

// --- 防止 unused import ---
var _ = json.Marshal
