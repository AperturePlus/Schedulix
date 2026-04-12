package model

import "time"

// ClusterMetrics 集群级指标快照。
// 每次采集生成一个不可变快照，存入环形缓冲区。
type ClusterMetrics struct {
	Timestamp           time.Time `json:"timestamp"`
	Version             int       `json:"version"`               // 递增版本号
	TotalTasks          int       `json:"total_tasks"`
	CompletedTasks      int       `json:"completed_tasks"`
	FailedTasks         int       `json:"failed_tasks"`
	AvgScheduleDelayMs  float64   `json:"avg_schedule_delay_ms"`
	ResourceUtilization float64   `json:"resource_utilization"` // [0.0, 1.0]
}

// NodeMetrics 节点级指标快照。
type NodeMetrics struct {
	Timestamp     time.Time `json:"timestamp"`
	NodeID        string    `json:"node_id"`
	CurrentLoad   float64   `json:"current_load"` // [0.0, 1.0]
	AssignedTasks int       `json:"assigned_tasks"`
	FaultCount    int       `json:"fault_count"`
	UptimeMs      int64     `json:"uptime_ms"`
}
