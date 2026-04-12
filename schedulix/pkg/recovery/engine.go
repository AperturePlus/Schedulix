package recovery

import (
	"errors"
	"fmt"
	"time"

	"schedulix/pkg/model"
	"schedulix/pkg/queue"
	"schedulix/pkg/simulator"
)

// ─── 错误定义 ───────────────────────────────────────────────

var (
	// ErrRecoveryNotReady 恢复引擎未正确初始化。
	ErrRecoveryNotReady = errors.New("recovery engine is not ready: missing cluster or queue")

	// ErrNodeNotFound 故障事件中的节点在集群中不存在。
	ErrNodeNotFound = errors.New("faulted node not found in cluster")

	// ErrTaskNotFound 任务在集群中找不到（可能已被其他恢复流程处理）。
	ErrTaskNotFound = errors.New("task not found")

	// ErrMaxMigrationsExceeded 任务迁移次数超过上限。
	ErrMaxMigrationsExceeded = errors.New("task exceeded maximum migration count")

	// ErrRequeueFailed 任务重新入队失败（队列满或其他错误）。
	ErrRequeueFailed = errors.New("failed to requeue task after fault")
)

// MaxMigrationCount 最大迁移次数。超过此值的任务标记为 Failed。
const MaxMigrationCount = 3

// RecoveryEngine 容灾恢复引擎。
//
// 鲁棒性设计：
//   - 单个任务的恢复失败不影响其他任务
//   - 重新入队失败时记录错误但不 panic
//   - 所有操作记录详细日志
//   - 幂等性：同一故障事件多次处理不会重复迁移
type RecoveryEngine struct {
	cluster      *model.Cluster
	queue        *queue.TaskQueue
	checkpoints  *CheckpointStore
	recoveryLogs []RecoveryLog
	// processedEvents 记录已处理的事件 ID，防止重复处理
	processedEvents map[string]bool
}

// RecoveryLog 恢复操作日志。
type RecoveryLog struct {
	Timestamp     time.Time
	FaultType     simulator.FaultType
	NodeID        string
	AffectedTasks int
	RecoveredTasks int   // 成功恢复的任务数
	FailedTasks   int    // 标记为 Failed 的任务数
	Errors        int    // 处理过程中的错误数
	RecoveryMs    int64  // 恢复耗时（毫秒）
	Detail        string
}

// NewRecoveryEngine 创建容灾恢复引擎。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - cluster == nil 或 q == nil → 仍然创建，OnFault 时检查并返回 ErrRecoveryNotReady
// - 初始化 CheckpointStore、recoveryLogs、processedEvents
func NewRecoveryEngine(cluster *model.Cluster, q *queue.TaskQueue) *RecoveryEngine {
	// TODO: 实现
	panic("not implemented")
}

// OnFault 实现 EventHandler 接口 — 处理故障事件。
//
// TODO(learner): 实现此方法
// 鲁棒性要求（完整的错误处理链）：
// 1. event == nil → 返回错误（不 panic）
// 2. 检查 cluster 和 queue 是否就绪
// 3. 幂等性检查：event.ID 是否已处理过 → 是则跳过
// 4. 获取故障节点 → 不存在则返回 ErrNodeNotFound（包装节点 ID）
// 5. 获取节点上的 AssignedTasks（复制一份，因为后续会修改原切片）
// 6. 对每个任务：
//    a. 从检查点恢复进度（检查点不存在 → 从 0 开始，不是错误）
//    b. 增加 MigrationCount
//    c. 如果 MigrationCount >= MaxMigrationCount → 标记为 Failed，记录
//    d. 否则 → 重新提交到队列
//       - 入队失败 → 记录错误，继续处理下一个任务（不中断）
// 7. 记录 RecoveryLog（包含成功/失败/错误计数）
// 8. 标记事件为已处理
//
// 关键原则：单个任务的恢复失败不能阻止其他任务的恢复。
func (re *RecoveryEngine) OnFault(event *simulator.FaultEvent) error {
	// TODO: 实现
	panic("not implemented")
}

// OnRecovery 实现 EventHandler 接口 — 处理恢复事件。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. event == nil → 返回错误
// 2. 幂等性检查
// 3. 获取节点 → 不存在则返回错误
// 4. 更新节点状态为 Idle
// 5. 记录日志
func (re *RecoveryEngine) OnRecovery(event *simulator.FaultEvent) error {
	// TODO: 实现
	panic("not implemented")
}

// GetRecoveryLogs 返回所有恢复操作日志的副本。
func (re *RecoveryEngine) GetRecoveryLogs() []RecoveryLog {
	result := make([]RecoveryLog, len(re.recoveryLogs))
	copy(result, re.recoveryLogs)
	return result
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
