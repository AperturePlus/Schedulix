package simulator

import (
	"errors"
	"fmt"
	"sync"

	"schedulix/pkg/model"
)

// ─── 错误定义 ───────────────────────────────────────────────

var (
	// ErrSimulatorNotReady 模拟器未正确初始化。
	ErrSimulatorNotReady = errors.New("event simulator is not ready: missing config or cluster")

	// ErrSimulationAborted 模拟被中断（handler 返回错误）。
	ErrSimulationAborted = errors.New("simulation aborted due to handler error")

	// ErrHandlerPanic handler 处理事件时 panic。
	ErrHandlerPanic = errors.New("event handler panicked")
)

// EventHandler 事件处理回调接口。
// 采用观察者模式：EventSimulator 生成事件后遍历 handler 列表调用。
//
// 鲁棒性契约：
//   - 实现者应快速返回，不做长时间阻塞操作
//   - 实现者返回 error 时，模拟器记录错误但继续运行（不中断模拟）
//   - 实现者不应 panic（但模拟器会 recover）
type EventHandler interface {
	// OnFault 处理故障事件（宕机、网络延迟、性能降级）。
	OnFault(event *FaultEvent) error
	// OnRecovery 处理恢复事件。
	OnRecovery(event *FaultEvent) error
}

// EventSimulator 事件模拟引擎。
//
// 鲁棒性设计：
//   - Handler panic 被 recover，不影响其他 handler 和后续事件
//   - 单个 handler 返回错误时记录日志但继续通知其他 handler
//   - 事件日志加锁保护，支持并发读取
//   - 配置不合法时使用默认配置（优雅降级）
type EventSimulator struct {
	config     *EventConfig
	cluster    *model.Cluster
	handlers   []EventHandler
	eventLog   []*FaultEvent // 按时间顺序的事件日志
	handlerMu  sync.RWMutex // 保护 handlers 列表
	logMu      sync.RWMutex // 保护 eventLog
	errorLog   []error       // 记录 handler 处理错误（不中断模拟）
}

// NewEventSimulator 创建事件模拟器。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - config == nil → 使用 DefaultEventConfig()（优雅降级）
// - cluster == nil → 仍然创建，RunStepMode 时检查并返回 ErrSimulatorNotReady
// - 初始化所有切片为非 nil 空切片
func NewEventSimulator(config *EventConfig, cluster *model.Cluster) *EventSimulator {
	// TODO: 实现
	panic("not implemented")
}

// RegisterHandler 注册事件处理器。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - handler == nil → 静默忽略（不注册 nil handler）
// - 使用 handlerMu 保护并发注册
func (es *EventSimulator) RegisterHandler(handler EventHandler) {
	// TODO: 实现
	panic("not implemented")
}

// RunStepMode 运行时间步进模式。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. cluster == nil → 返回 ErrSimulatorNotReady
// 2. config.Validate() 失败 → 调用 config.Clamp() 自动修复，记录警告，继续运行
// 3. 每个时间步：
//    a. 遍历集群节点
//    b. 对每个节点做伯努利试验
//    c. 生成事件后，调用 notifyHandlers
// 4. notifyHandlers 内部：
//    a. 遍历所有 handler
//    b. 每个 handler 调用包裹在 defer recover() 中
//    c. handler 返回 error → 记录到 errorLog，继续下一个 handler
//    d. handler panic → recover，记录 ErrHandlerPanic，继续
// 5. 所有事件记录到 eventLog（加锁）
func (es *EventSimulator) RunStepMode() error {
	// TODO: 实现
	panic("not implemented")
}

// notifyHandlers 安全地通知所有 handler。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - 每个 handler 调用包裹在 safeCall 中
// - safeCall 使用 defer recover() 捕获 panic
// - 单个 handler 失败不影响其他 handler
func (es *EventSimulator) notifyHandlers(event *FaultEvent) {
	// TODO: 实现
	panic("not implemented")
}

// GetEventLog 返回事件日志的副本（线程安全）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - 返回副本而非原始切片（防止调用方修改内部状态）
// - 使用 logMu.RLock() 保护
func (es *EventSimulator) GetEventLog() []*FaultEvent {
	// TODO: 实现
	panic("not implemented")
}

// GetErrors 返回模拟过程中 handler 产生的错误列表。
func (es *EventSimulator) GetErrors() []error {
	es.logMu.RLock()
	defer es.logMu.RUnlock()
	result := make([]error, len(es.errorLog))
	copy(result, es.errorLog)
	return result
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
