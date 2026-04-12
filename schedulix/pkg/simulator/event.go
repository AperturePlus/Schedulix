package simulator

import "time"

// FaultType 故障类型枚举。
type FaultType int

const (
	FaultNodeDown     FaultType = iota // 节点宕机：状态变为 Offline，所有任务需迁移
	FaultNetworkDelay                  // 网络延迟：模拟调度延迟增大
	FaultDegraded                      // 性能降级：算力折扣（如减半）
	FaultRecovery                      // 节点恢复：状态恢复为 Idle
)

// FaultEvent 模拟故障事件。
// 事件是不可变的：一旦创建不应修改。
type FaultEvent struct {
	ID        string    `json:"id"`
	Type      FaultType `json:"type"`
	NodeID    string    `json:"node_id"`
	Timestamp time.Time `json:"timestamp"`
	Detail    string    `json:"detail,omitempty"`
}
