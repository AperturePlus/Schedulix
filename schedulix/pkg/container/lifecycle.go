package container

import (
	"errors"
	"fmt"

	"schedulix/pkg/model"
)

// ─── 错误定义 ───────────────────────────────────────────────

var (
	// ErrInvalidStateTransition 非法状态转换错误。
	ErrInvalidStateTransition = errors.New("invalid container state transition")

	// ErrNilContainer 容器指针为 nil。
	ErrNilContainer = errors.New("container is nil")

	// ErrContainerDestroyed 容器已销毁，不能再操作。
	ErrContainerDestroyed = errors.New("container is already destroyed")
)

// ContainerLifecycle 容器状态变更的事件订阅接口（观察者模式）。
//
// 鲁棒性契约：
//   - 实现者不应 panic（调用方会 recover）
//   - 实现者应快速返回
//   - 回调不返回 error（状态变更已发生，不可阻止）
type ContainerLifecycle interface {
	OnStateChange(containerID string, oldState, newState model.ContainerState)
}

// validTransitions 合法的状态转换表。
var validTransitions = map[model.ContainerState][]model.ContainerState{
	model.ContainerCreated: {model.ContainerRunning},
	model.ContainerRunning: {model.ContainerStopped},
	model.ContainerStopped: {model.ContainerDestroyed},
}

// IsValidTransition 检查状态转换是否合法。
//
// TODO(learner): 实现此方法
// 提示：在 validTransitions 中查找 from 状态，检查 to 是否在允许列表中
func IsValidTransition(from, to model.ContainerState) bool {
	// TODO: 实现
	panic("not implemented")
}

// TransitionState 执行容器状态转换。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. c == nil → 返回 ErrNilContainer
// 2. c.State == ContainerDestroyed → 返回 ErrContainerDestroyed（已销毁不可操作）
// 3. 检查转换是否合法（IsValidTransition）
//    - 不合法 → 返回 fmt.Errorf("%w: %v → %v", ErrInvalidStateTransition, from, to)
// 4. 合法 → 更新容器状态
// 5. 安全通知所有观察者：
//    - 跳过 nil 观察者
//    - 每个观察者调用包裹在 defer recover() 中
//    - 单个观察者 panic 不影响其他观察者
func TransitionState(c *model.Container, newState model.ContainerState, observers []ContainerLifecycle) error {
	// TODO: 实现
	panic("not implemented")
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
