package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"schedulix/pkg/model"
)

// ============================================================
// 状态机测试
// ============================================================

func TestIsValidTransition(t *testing.T) {
	tests := []struct {
		from  model.ContainerState
		to    model.ContainerState
		valid bool
	}{
		{model.ContainerCreated, model.ContainerRunning, true},
		{model.ContainerRunning, model.ContainerStopped, true},
		{model.ContainerStopped, model.ContainerDestroyed, true},
		// 非法转换
		{model.ContainerCreated, model.ContainerStopped, false},
		{model.ContainerCreated, model.ContainerDestroyed, false},
		{model.ContainerRunning, model.ContainerCreated, false},
		{model.ContainerStopped, model.ContainerRunning, false},
		{model.ContainerDestroyed, model.ContainerCreated, false},
	}

	for _, tt := range tests {
		got := IsValidTransition(tt.from, tt.to)
		assert.Equal(t, tt.valid, got, "%v → %v", tt.from, tt.to)
	}
}

func TestTransitionState_NilContainer(t *testing.T) {
	err := TransitionState(nil, model.ContainerRunning, nil)
	assert.ErrorIs(t, err, ErrNilContainer)
}

func TestTransitionState_InvalidTransition(t *testing.T) {
	c := &model.Container{ID: "c-1", State: model.ContainerCreated}
	err := TransitionState(c, model.ContainerDestroyed, nil)
	assert.ErrorIs(t, err, ErrInvalidStateTransition)
}

func TestTransitionState_FullLifecycle(t *testing.T) {
	c := &model.Container{ID: "c-1", State: model.ContainerCreated}

	require.NoError(t, TransitionState(c, model.ContainerRunning, nil))
	assert.Equal(t, model.ContainerRunning, c.State)

	require.NoError(t, TransitionState(c, model.ContainerStopped, nil))
	assert.Equal(t, model.ContainerStopped, c.State)

	require.NoError(t, TransitionState(c, model.ContainerDestroyed, nil))
	assert.Equal(t, model.ContainerDestroyed, c.State)
}

// ============================================================
// 观察者测试
// ============================================================

type mockObserver struct {
	calls []struct{ containerID, from, to string }
}

func (m *mockObserver) OnStateChange(containerID string, oldState, newState model.ContainerState) {
	m.calls = append(m.calls, struct{ containerID, from, to string }{
		containerID, oldState.String(), newState.String(),
	})
}

func (s model.ContainerState) String() string {
	names := map[model.ContainerState]string{
		model.ContainerCreated:   "created",
		model.ContainerRunning:   "running",
		model.ContainerStopped:   "stopped",
		model.ContainerDestroyed: "destroyed",
	}
	return names[s]
}

func TestTransitionState_NotifiesObservers(t *testing.T) {
	obs := &mockObserver{}
	c := &model.Container{ID: "c-1", State: model.ContainerCreated}

	TransitionState(c, model.ContainerRunning, []ContainerLifecycle{obs})

	require.Len(t, obs.calls, 1)
	assert.Equal(t, "c-1", obs.calls[0].containerID)
	assert.Equal(t, "created", obs.calls[0].from)
	assert.Equal(t, "running", obs.calls[0].to)
}

func TestTransitionState_ObserverPanicIsolation(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建一个会 panic 的观察者
	// 2. 创建一个正常的观察者
	// 3. 执行状态转换
	// 4. 验证正常观察者仍然收到通知
	// 5. 验证状态转换成功完成
}

// ============================================================
// ContainerRuntime 测试
// ============================================================

func TestCreateContainer(t *testing.T) {
	cr := NewContainerRuntime()
	node := &model.GPU_Node{ID: "n-1", MemoryTotal: 8000}

	c, err := cr.CreateContainer("c-1", node, 100, 2000)
	require.NoError(t, err)
	assert.Equal(t, "c-1", c.ID)
	assert.Equal(t, model.ContainerCreated, c.State)
}

func TestCreateContainer_InsufficientResources(t *testing.T) {
	cr := NewContainerRuntime()
	node := &model.GPU_Node{ID: "n-1", MemoryTotal: 1000}

	_, err := cr.CreateContainer("c-1", node, 100, 2000)
	assert.ErrorIs(t, err, ErrInsufficientResources)
}

func TestCreateContainer_DuplicateID(t *testing.T) {
	cr := NewContainerRuntime()
	node := &model.GPU_Node{ID: "n-1", MemoryTotal: 8000}

	cr.CreateContainer("c-1", node, 100, 1000)
	_, err := cr.CreateContainer("c-1", node, 100, 1000)
	assert.ErrorIs(t, err, ErrDuplicateContainerID)
}

func TestCreateContainer_NilNode(t *testing.T) {
	cr := NewContainerRuntime()
	_, err := cr.CreateContainer("c-1", nil, 100, 1000)
	assert.ErrorIs(t, err, ErrNilHostNode)
}

func TestCreateContainer_MultipleOnSameNode(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建节点（8000MB）
	// 2. 创建容器 A（3000MB）→ 成功
	// 3. 创建容器 B（3000MB）→ 成功
	// 4. 创建容器 C（3000MB）→ 失败（3000+3000+3000 > 8000）
}

func TestContainerRuntime_FullLifecycle(t *testing.T) {
	// TODO(learner): 实现
	// Create → Start → Stop → Destroy
	// 验证每步状态正确
}
