package orchestrator

import (
	"time"

	"schedulix/pkg/model"
)

// ─── K8s 核心概念映射 ──────────────────────────────────────
//
// 学习要点：
//   Kubernetes 的核心抽象层次：
//
//   Pod        — 最小调度单元，包含一个或多个容器
//   ReplicaSet — 确保指定数量的 Pod 副本始终运行
//   Deployment — 管理 ReplicaSet，支持滚动更新和回滚
//   Service    — 为一组 Pod 提供稳定的访问入口（服务发现）
//   Node       — 运行 Pod 的工作节点（对应我们的 GPU_Node）
//
//   Schedulix 中的映射：
//   K8s Pod        → Pod（包含一个 Task + 资源配额）
//   K8s ReplicaSet → ReplicaSet（维护 N 个 Pod 副本）
//   K8s Deployment → Deployment（滚动更新 ReplicaSet）
//   K8s Service    → Service（负载均衡到一组 Pod）
//   K8s Node       → GPU_Node（已有）

// ─── Pod ────────────────────────────────────────────────────

// PodPhase Pod 生命周期阶段（对应 K8s Pod Phase）。
type PodPhase int

const (
	PodPending   PodPhase = iota // 等待调度
	PodRunning                   // 已调度，容器运行中
	PodSucceeded                 // 所有容器正常退出
	PodFailed                    // 至少一个容器异常退出
	PodUnknown                   // 无法获取状态（节点失联）
)

// RestartPolicy Pod 重启策略。
type RestartPolicy int

const (
	RestartAlways    RestartPolicy = iota // 总是重启（默认）
	RestartOnFailure                      // 仅失败时重启
	RestartNever                          // 从不重启
)

// Pod 最小调度单元。
//
// K8s 中 Pod 可以包含多个容器（sidecar 模式），
// Schedulix 简化为一个 Pod 对应一个 Task。
type Pod struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Namespace     string              `json:"namespace"`      // 命名空间隔离
	Phase         PodPhase            `json:"phase"`
	NodeID        string              `json:"node_id"`        // 调度到的节点
	Resource      model.ResourceRequirement `json:"resource"`
	RestartPolicy RestartPolicy       `json:"restart_policy"`
	RestartCount  int                 `json:"restart_count"`
	Labels        map[string]string   `json:"labels"`         // 标签选择器用
	CreatedAt     time.Time           `json:"created_at"`
	StartedAt     *time.Time          `json:"started_at,omitempty"`
	FinishedAt    *time.Time          `json:"finished_at,omitempty"`
	// 健康检查
	LivenessProbe  *Probe             `json:"liveness_probe,omitempty"`
	ReadinessProbe *Probe             `json:"readiness_probe,omitempty"`
}

// Probe 健康探针配置。
//
// K8s 中有三种探针：
//   - Liveness: 容器是否存活？失败 → 重启容器
//   - Readiness: 容器是否就绪？失败 → 从 Service 中移除
//   - Startup: 容器是否启动完成？（Schedulix 暂不实现）
type Probe struct {
	IntervalMs     int64 `json:"interval_ms"`      // 检查间隔
	TimeoutMs      int64 `json:"timeout_ms"`        // 超时时间
	FailureThreshold int `json:"failure_threshold"` // 连续失败多少次判定为不健康
	SuccessThreshold int `json:"success_threshold"` // 连续成功多少次判定为健康
}

// ─── ReplicaSet ─────────────────────────────────────────────

// ReplicaSet 副本集，确保指定数量的 Pod 始终运行。
//
// 核心行为（控制循环 / Reconciliation Loop）：
//   期望状态：replicas = 3
//   当前状态：2 个 Pod 运行中
//   动作：创建 1 个新 Pod
//
//   期望状态：replicas = 3
//   当前状态：4 个 Pod 运行中
//   动作：终止 1 个 Pod
type ReplicaSet struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Replicas  int               `json:"replicas"`       // 期望副本数
	Selector  map[string]string `json:"selector"`       // 标签选择器
	Template  PodTemplate       `json:"template"`       // Pod 模板
	// 状态
	ReadyReplicas   int `json:"ready_replicas"`
	CurrentReplicas int `json:"current_replicas"`
}

// PodTemplate Pod 创建模板。
type PodTemplate struct {
	Labels        map[string]string         `json:"labels"`
	Resource      model.ResourceRequirement `json:"resource"`
	RestartPolicy RestartPolicy             `json:"restart_policy"`
	LivenessProbe  *Probe                   `json:"liveness_probe,omitempty"`
	ReadinessProbe *Probe                   `json:"readiness_probe,omitempty"`
}

// ─── Deployment ─────────────────────────────────────────────

// DeploymentStrategy 部署策略。
type DeploymentStrategy int

const (
	// RollingUpdate 滚动更新：逐步替换旧 Pod 为新 Pod。
	// 保证在更新过程中始终有一定数量的 Pod 可用。
	RollingUpdate DeploymentStrategy = iota

	// Recreate 重建：先删除所有旧 Pod，再创建新 Pod。
	// 简单但有停机时间。
	Recreate
)

// RollingUpdateConfig 滚动更新配置。
type RollingUpdateConfig struct {
	MaxUnavailable int `json:"max_unavailable"` // 更新期间最多不可用的 Pod 数
	MaxSurge       int `json:"max_surge"`       // 更新期间最多额外创建的 Pod 数
}

// Deployment 部署，管理 ReplicaSet 的版本和更新。
type Deployment struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Replicas  int                 `json:"replicas"`
	Strategy  DeploymentStrategy  `json:"strategy"`
	RollingUpdate *RollingUpdateConfig `json:"rolling_update,omitempty"`
	Template  PodTemplate         `json:"template"`
	// 版本管理
	Revision       int            `json:"revision"`
	ReplicaSets    []string       `json:"replica_set_ids"` // 历史 ReplicaSet ID
	MaxRevisions   int            `json:"max_revisions"`   // 保留的历史版本数（用于回滚）
}

// ─── Service ────────────────────────────────────────────────

// ServiceType 服务类型。
type ServiceType int

const (
	ClusterIP ServiceType = iota // 集群内部访问（默认）
	NodePort                     // 通过节点端口暴露
)

// Service 服务发现与负载均衡。
//
// Service 通过标签选择器（Selector）找到匹配的 Pod，
// 将流量负载均衡到这些 Pod 上。
// Pod 的 IP 会变（重启后新 IP），但 Service 的名称不变 → 稳定的访问入口。
type Service struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Type      ServiceType       `json:"type"`
	Selector  map[string]string `json:"selector"`  // 匹配 Pod 的标签
	Port      int               `json:"port"`
	TargetPort int              `json:"target_port"`
	// 运行时状态
	Endpoints []string          `json:"endpoints"` // 匹配的 Pod ID 列表
}

// ─── Namespace ──────────────────────────────────────────────

// Namespace 命名空间，用于资源隔离。
//
// 不同命名空间中的资源互不可见（除非显式跨命名空间访问）。
// 常见用法：dev / staging / production 环境隔离。
type Namespace struct {
	Name         string            `json:"name"`
	Labels       map[string]string `json:"labels,omitempty"`
	ResourceQuota *ResourceQuota   `json:"resource_quota,omitempty"`
}

// ResourceQuota 命名空间级别的资源配额限制。
type ResourceQuota struct {
	MaxPods         int `json:"max_pods"`
	MaxCPU          int `json:"max_cpu"`
	MaxMemory       int `json:"max_memory"`        // MB
}
