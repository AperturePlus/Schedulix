package container

import (
	"errors"
	"fmt"
	"sync"

	"schedulix/pkg/model"
)

// ─── 错误定义 ───────────────────────────────────────────────

var (
	// ErrInsufficientResources 资源不足错误。
	ErrInsufficientResources = errors.New("insufficient resources on host node")

	// ErrContainerNotFound 容器不存在。
	ErrContainerNotFound = errors.New("container not found")

	// ErrDuplicateContainerID 容器 ID 重复。
	ErrDuplicateContainerID = errors.New("container ID already exists")

	// ErrInvalidContainerID 容器 ID 为空。
	ErrInvalidContainerID = errors.New("container ID must be non-empty")

	// ErrNilHostNode 宿主节点为 nil。
	ErrNilHostNode = errors.New("host node is nil")

	// ErrInvalidResourceQuota 资源配额不合法。
	ErrInvalidResourceQuota = errors.New("resource quota must be positive")
)

// ContainerRuntime 容器运行时，管理容器的创建和资源分配。
//
// 鲁棒性设计：
//   - 线程安全（sync.RWMutex）
//   - 所有输入参数验证
//   - 容器 ID 唯一性检查
//   - 资源配额验证
//   - 观察者 panic 不影响运行时
type ContainerRuntime struct {
	containers map[string]*model.Container
	observers  []ContainerLifecycle
	mu         sync.RWMutex
}

// NewContainerRuntime 创建容器运行时。
func NewContainerRuntime() *ContainerRuntime {
	return &ContainerRuntime{
		containers: make(map[string]*model.Container),
	}
}

// RegisterObserver 注册容器生命周期观察者。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - obs == nil → 静默忽略
// - 使用 mu.Lock() 保护
func (cr *ContainerRuntime) RegisterObserver(obs ContainerLifecycle) {
	// TODO: 实现
	panic("not implemented")
}

// CreateContainer 在指定节点上创建容器。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. id == "" → 返回 nil, ErrInvalidContainerID
// 2. hostNode == nil → 返回 nil, ErrNilHostNode
// 3. cpuShares <= 0 或 memoryLimit <= 0 → 返回 nil, ErrInvalidResourceQuota
// 4. id 已存在 → 返回 nil, ErrDuplicateContainerID（幂等性：不重复创建）
// 5. 检查宿主节点剩余资源：
//    - 计算节点上已有容器的内存总和（遍历 containers，筛选同一 hostNode）
//    - 加上新容器的 memoryLimit，不能超过 hostNode.MemoryTotal
//    - 超过 → 返回 nil, ErrInsufficientResources
// 6. 创建容器（状态为 Created），加入 containers map
// 7. 使用 mu.Lock() 保护整个操作
func (cr *ContainerRuntime) CreateContainer(id string, hostNode *model.GPU_Node, cpuShares, memoryLimit int) (*model.Container, error) {
	// TODO: 实现
	panic("not implemented")
}

// StartContainer 启动容器（Created → Running）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. containerID == "" → 返回 ErrInvalidContainerID
// 2. 容器不存在 → 返回 ErrContainerNotFound
// 3. 调用 TransitionState（内部处理非法转换和观察者通知）
func (cr *ContainerRuntime) StartContainer(containerID string) error {
	// TODO: 实现
	panic("not implemented")
}

// StopContainer 停止容器（Running → Stopped）。
//
// TODO(learner): 实现此方法（同 StartContainer 的鲁棒性要求）
func (cr *ContainerRuntime) StopContainer(containerID string) error {
	// TODO: 实现
	panic("not implemented")
}

// DestroyContainer 销毁容器（Stopped → Destroyed）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 同 StartContainer 的输入验证
// 2. 状态转换成功后，从 containers map 中移除（释放资源）
//    — 或者保留但标记为 Destroyed（取决于你的设计选择，两种都合理）
func (cr *ContainerRuntime) DestroyContainer(containerID string) error {
	// TODO: 实现
	panic("not implemented")
}

// GetContainer 获取容器信息（线程安全）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - 使用 mu.RLock() 保护
// - 返回副本而非指针（防止调用方修改内部状态）— 或者返回指针但文档说明只读
func (cr *ContainerRuntime) GetContainer(containerID string) (*model.Container, bool) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	c, ok := cr.containers[containerID]
	return c, ok
}

// ListContainers 列出所有容器（线程安全）。
func (cr *ContainerRuntime) ListContainers() []*model.Container {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	result := make([]*model.Container, 0, len(cr.containers))
	for _, c := range cr.containers {
		result = append(result, c)
	}
	return result
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
