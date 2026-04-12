package model

import (
	"encoding/json"
	"fmt"
)

// NodeStatus 节点状态枚举。
// 使用 int + iota：编译期类型安全，switch 可检查穷举，比较 O(1)。
type NodeStatus int

const (
	NodeStatusIdle     NodeStatus = iota // 空闲，可接受新任务
	NodeStatusBusy                       // 忙碌，资源已分配但未耗尽
	NodeStatusOffline                    // 离线/宕机，不可调度
	NodeStatusDegraded                   // 性能降级，算力折扣运行
)

// nodeStatusNames 状态名称映射，用于 JSON 序列化。
var nodeStatusNames = map[NodeStatus]string{
	NodeStatusIdle:     "idle",
	NodeStatusBusy:     "busy",
	NodeStatusOffline:  "offline",
	NodeStatusDegraded: "degraded",
}

// nodeStatusValues 反向映射，用于 JSON 反序列化。
var nodeStatusValues = map[string]NodeStatus{
	"idle":     NodeStatusIdle,
	"busy":     NodeStatusBusy,
	"offline":  NodeStatusOffline,
	"degraded": NodeStatusDegraded,
}

// String 返回状态的字符串表示。
func (s NodeStatus) String() string {
	if name, ok := nodeStatusNames[s]; ok {
		return name
	}
	return fmt.Sprintf("unknown(%d)", int(s))
}

// MarshalJSON 将 NodeStatus 序列化为 JSON 字符串。
//
// TODO(learner): 实现此方法
// 提示：使用 json.Marshal 将 nodeStatusNames 中对应的字符串序列化
// 如果状态值不在映射中，返回错误
func (s NodeStatus) MarshalJSON() ([]byte, error) {
	// TODO: 实现 JSON 序列化
	panic("not implemented")
}

// UnmarshalJSON 从 JSON 字符串反序列化为 NodeStatus。
//
// TODO(learner): 实现此方法
// 提示：
// 1. 先将 data 反序列化为 string
// 2. 在 nodeStatusValues 中查找对应的 NodeStatus
// 3. 找不到则返回错误
func (s *NodeStatus) UnmarshalJSON(data []byte) error {
	// TODO: 实现 JSON 反序列化
	panic("not implemented")
}

// --- TaskStatus ---

// TaskStatus 任务状态。
// 状态转换规则：
//
//	Pending → Running（被调度到节点）
//	Running → Completed（正常完成）
//	Running → Migrating（节点故障，迁移中）
//	Migrating → Pending（迁移后重新入队）
//	Migrating → Failed（三次迁移失败）
//	Pending → Failed（超时未调度）
type TaskStatus int

const (
	TaskStatusPending   TaskStatus = iota // 等待调度
	TaskStatusRunning                     // 执行中
	TaskStatusCompleted                   // 已完成
	TaskStatusFailed                      // 失败（不可恢复）
	TaskStatusMigrating                   // 迁移中
)

var taskStatusNames = map[TaskStatus]string{
	TaskStatusPending:   "pending",
	TaskStatusRunning:   "running",
	TaskStatusCompleted: "completed",
	TaskStatusFailed:    "failed",
	TaskStatusMigrating: "migrating",
}

var taskStatusValues = map[string]TaskStatus{
	"pending":   TaskStatusPending,
	"running":   TaskStatusRunning,
	"completed": TaskStatusCompleted,
	"failed":    TaskStatusFailed,
	"migrating": TaskStatusMigrating,
}

func (s TaskStatus) String() string {
	if name, ok := taskStatusNames[s]; ok {
		return name
	}
	return fmt.Sprintf("unknown(%d)", int(s))
}

// MarshalJSON 将 TaskStatus 序列化为 JSON 字符串。
// TODO(learner): 参考 NodeStatus.MarshalJSON 实现
func (s TaskStatus) MarshalJSON() ([]byte, error) {
	panic("not implemented")
}

// UnmarshalJSON 从 JSON 字符串反序列化为 TaskStatus。
// TODO(learner): 参考 NodeStatus.UnmarshalJSON 实现
func (s *TaskStatus) UnmarshalJSON(data []byte) error {
	panic("not implemented")
}

// --- 以下为占位，防止编译器报 unused import ---
var _ json.Marshaler = NodeStatus(0)
var _ json.Unmarshaler = (*NodeStatus)(nil)
