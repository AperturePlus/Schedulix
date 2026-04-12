package model

import "time"

// Checkpoint 任务执行检查点。
// 存储在内存 map[string]*Checkpoint 中（key 为 TaskID）。
// 每次保存覆盖前一个检查点（每个任务只保留最新一个）。
type Checkpoint struct {
	TaskID    string    `json:"task_id"`
	Progress  float64   `json:"progress"`  // 保存时的执行进度 [0.0, 1.0]
	Timestamp time.Time `json:"timestamp"`
	NodeID    string    `json:"node_id"` // 保存时所在节点
}
