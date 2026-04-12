package model

// ContainerState 容器状态。
// 合法状态转换：Created→Running, Running→Stopped, Stopped→Destroyed
// 非法转换返回 ErrInvalidStateTransition。
type ContainerState int

const (
	ContainerCreated   ContainerState = iota // 已创建
	ContainerRunning                         // 运行中
	ContainerStopped                         // 已停止
	ContainerDestroyed                       // 已销毁
)

// Container 模拟容器。
type Container struct {
	ID          string         `json:"id"`
	State       ContainerState `json:"state"`
	HostNodeID  string         `json:"host_node_id"`
	CPUShares   int            `json:"cpu_shares"`   // CPU 份额
	MemoryLimit int            `json:"memory_limit"` // 内存限制（MB）
	TaskID      string         `json:"task_id,omitempty"`
}
