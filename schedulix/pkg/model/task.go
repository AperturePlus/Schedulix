package model

import "time"

// ResourceRequirement 资源需求。
// 独立结构体便于扩展（未来可增加 GPU 显存、网络带宽等维度）。
type ResourceRequirement struct {
	ComputePower int `json:"compute_power"` // 所需算力（TFLOPS）
	Memory       int `json:"memory"`        // 所需内存（MB）
}

// Task 用户提交的计算任务。
type Task struct {
	ID              string              `json:"id"`
	Resource        ResourceRequirement `json:"resource"`
	Priority        int                 `json:"priority"`          // 数值越大优先级越高
	EstimatedTimeMs int64               `json:"estimated_time_ms"` // 预计执行时间（毫秒）
	SubmitTime      time.Time           `json:"submit_time"`
	Status          TaskStatus          `json:"status"`
	AssignedNodeID  string              `json:"assigned_node_id,omitempty"`
	MigrationCount  int                 `json:"migration_count"` // 迁移次数，>=3 时标记失败
	Progress        float64             `json:"progress"`        // 执行进度 [0.0, 1.0]
}
