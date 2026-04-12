package orchestrator

import (
	"errors"
	"fmt"
	"sync"
)

// ─── 错误定义 ───────────────────────────────────────────────

var (
	ErrDeploymentNotFound = errors.New("deployment not found")
	ErrNoRevisionToRollback = errors.New("no previous revision to rollback to")
	ErrRolloutInProgress = errors.New("a rollout is already in progress")
)

// DeploymentController 部署控制器。
//
// 学习要点 — 滚动更新（Rolling Update）：
//
//   假设当前有 3 个 v1 Pod，要更新到 v2：
//
//   MaxUnavailable=1, MaxSurge=1 的滚动更新过程：
//
//   步骤 1: 创建 1 个 v2 Pod（surge）
//     v1: [●] [●] [●]
//     v2: [○]              ← 新建
//     总数: 4（3 + surge 1）
//
//   步骤 2: v2 Pod 就绪后，终止 1 个 v1 Pod
//     v1: [●] [●] [✗]     ← 终止
//     v2: [●]              ← 就绪
//     总数: 3
//
//   步骤 3: 创建 1 个 v2 Pod
//     v1: [●] [●]
//     v2: [●] [○]          ← 新建
//     总数: 4
//
//   步骤 4: v2 就绪，终止 1 个 v1
//     v1: [●] [✗]
//     v2: [●] [●]
//     总数: 3
//
//   步骤 5-6: 重复...
//     v1: (空)
//     v2: [●] [●] [●]
//     完成！
//
//   关键保证：
//   - 任何时刻可用 Pod 数 >= replicas - maxUnavailable
//   - 任何时刻总 Pod 数 <= replicas + maxSurge
type DeploymentController struct {
	rsController *ReplicaSetController
	deployments  map[string]*Deployment
	mu           sync.RWMutex
}

// NewDeploymentController 创建部署控制器。
//
// TODO(learner): 实现此方法
func NewDeploymentController(rsController *ReplicaSetController) *DeploymentController {
	// TODO: 实现
	panic("not implemented")
}

// CreateDeployment 创建部署。
//
// TODO(learner): 实现此方法
// 步骤：
// 1. 验证输入
// 2. 创建初始 ReplicaSet（revision = 1）
// 3. 存储 Deployment
func (dc *DeploymentController) CreateDeployment(dep *Deployment) error {
	// TODO: 实现
	panic("not implemented")
}

// UpdateDeployment 更新部署（触发滚动更新或重建）。
//
// TODO(learner): 实现此方法
// 步骤：
// 1. 根据 dep.Strategy 选择更新方式
// 2. Recreate 策略：
//    a. 将旧 ReplicaSet 缩容到 0
//    b. 创建新 ReplicaSet
// 3. RollingUpdate 策略：
//    a. 创建新 ReplicaSet（replicas = 0）
//    b. 逐步扩容新 RS、缩容旧 RS
//    c. 每步遵守 MaxUnavailable 和 MaxSurge 约束
// 4. 递增 revision
// 5. 清理超过 MaxRevisions 的旧 ReplicaSet
//
// 鲁棒性要求：
// - 更新过程中任何步骤失败 → 停止更新，保持当前状态（不回滚）
// - 用户可以手动 Rollback 到上一个版本
func (dc *DeploymentController) UpdateDeployment(depID string, newTemplate PodTemplate) error {
	// TODO: 实现
	panic("not implemented")
}

// RollingUpdate 执行滚动更新的核心逻辑。
//
// TODO(learner): 实现此方法
// 这是最复杂的方法，需要仔细处理并发和状态一致性。
//
// 算法（每一步）：
// 1. 计算当前可用 Pod 数 = oldRS.ReadyReplicas + newRS.ReadyReplicas
// 2. 计算可以终止的旧 Pod 数 = min(oldRS.CurrentReplicas, maxUnavailable - (desired - available))
// 3. 计算可以创建的新 Pod 数 = min(desired - newRS.CurrentReplicas, maxSurge - (total - desired))
// 4. 缩容旧 RS
// 5. 扩容新 RS
// 6. 等待新 Pod 就绪
// 7. 重复直到旧 RS 缩容到 0
func (dc *DeploymentController) rollingUpdate(dep *Deployment, oldRSID, newRSID string) error {
	// TODO: 实现
	panic("not implemented")
}

// Rollback 回滚到上一个版本。
//
// TODO(learner): 实现此方法
// 步骤：
// 1. 找到上一个 ReplicaSet（revision - 1）
// 2. 如果不存在 → ErrNoRevisionToRollback
// 3. 将当前 RS 缩容到 0
// 4. 将上一个 RS 扩容到 desired replicas
// 5. 递增 revision（回滚也是一个新版本）
func (dc *DeploymentController) Rollback(depID string) error {
	// TODO: 实现
	panic("not implemented")
}

// GetDeploymentStatus 获取部署状态摘要。
//
// TODO(learner): 实现此方法
func (dc *DeploymentController) GetDeploymentStatus(depID string) (*DeploymentStatus, error) {
	// TODO: 实现
	panic("not implemented")
}

// DeploymentStatus 部署状态摘要。
type DeploymentStatus struct {
	Name              string `json:"name"`
	DesiredReplicas   int    `json:"desired_replicas"`
	CurrentReplicas   int    `json:"current_replicas"`
	ReadyReplicas     int    `json:"ready_replicas"`
	AvailableReplicas int    `json:"available_replicas"`
	Revision          int    `json:"revision"`
	UpdateInProgress  bool   `json:"update_in_progress"`
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
