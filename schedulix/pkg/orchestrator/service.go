package orchestrator

import (
	"errors"
	"fmt"
	"sync"
)

// ─── 错误定义 ───────────────────────────────────────────────

var (
	ErrServiceNotFound  = errors.New("service not found")
	ErrNoEndpoints      = errors.New("service has no ready endpoints")
	ErrDuplicateService = errors.New("service name already exists in namespace")
)

// ServiceController 服务控制器 — 服务发现与负载均衡。
//
// 学习要点：
//   为什么需要 Service？
//   - Pod 的 IP 是临时的（重启后变化）
//   - 直接访问 Pod IP 不可靠
//   - Service 提供稳定的名称（如 "my-app.default"）
//   - Service 自动将流量分发到匹配的 Pod
//
//   工作原理：
//   1. Service 定义标签选择器（如 app=web）
//   2. 控制器持续监控匹配的 Pod
//   3. 只有 Ready 状态的 Pod 才加入 Endpoints
//   4. 请求到达 Service → 负载均衡到某个 Endpoint
//
//   DNS 解析（简化版）：
//   "my-service.my-namespace" → 解析到匹配的 Pod 列表
type ServiceController struct {
	services map[string]*Service // serviceID → Service
	pods     map[string]*Pod     // 引用，用于标签匹配
	mu       sync.RWMutex
}

// NewServiceController 创建服务控制器。
//
// TODO(learner): 实现此方法
func NewServiceController() *ServiceController {
	// TODO: 实现
	panic("not implemented")
}

// CreateService 创建服务。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 验证 svc 不为 nil
// 2. 验证 Name 和 Namespace 非空
// 3. 检查同 Namespace 下是否已有同名 Service → ErrDuplicateService
// 4. 初始化 Endpoints 为空切片
// 5. 调用 reconcileEndpoints 填充初始端点
func (sc *ServiceController) CreateService(svc *Service) error {
	// TODO: 实现
	panic("not implemented")
}

// Resolve 服务发现 — 根据服务名解析到一个可用的 Pod ID。
//
// TODO(learner): 实现此方法
// 这模拟了 K8s 中 DNS 解析 + kube-proxy 负载均衡的过程。
//
// 步骤：
// 1. 查找 Service
// 2. 如果 Endpoints 为空 → ErrNoEndpoints
// 3. 从 Endpoints 中选择一个（轮询或随机）
// 4. 返回选中的 Pod ID
//
// 鲁棒性要求：
// - 选中的 Pod 可能已经不存在（刚被终止）→ 重试下一个
// - 所有 Endpoint 都不可用 → ErrNoEndpoints
func (sc *ServiceController) Resolve(serviceName, namespace string) (string, error) {
	// TODO: 实现
	panic("not implemented")
}

// ReconcileEndpoints 更新服务的端点列表。
//
// TODO(learner): 实现此方法
// 当 Pod 状态变化时调用（创建、就绪、终止）。
//
// 步骤：
// 1. 获取 Service 的标签选择器
// 2. 遍历所有 Pod，找到标签匹配且 Phase == PodRunning 的 Pod
// 3. 更新 Service.Endpoints
//
// 鲁棒性要求：
// - Pod 列表可能在遍历过程中变化（并发修改）→ 使用快照
func (sc *ServiceController) ReconcileEndpoints(serviceID string, allPods map[string]*Pod) error {
	// TODO: 实现
	panic("not implemented")
}

// SetPods 更新 Pod 引用（供 ReconcileEndpoints 使用）。
func (sc *ServiceController) SetPods(pods map[string]*Pod) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.pods = pods
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
