package orchestrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"schedulix/pkg/model"
)

// ============================================================
// 标签匹配测试
// ============================================================

func TestMatchLabels(t *testing.T) {
	tests := []struct {
		name     string
		pod      map[string]string
		selector map[string]string
		want     bool
	}{
		{"exact match", map[string]string{"app": "web"}, map[string]string{"app": "web"}, true},
		{"superset matches", map[string]string{"app": "web", "env": "prod"}, map[string]string{"app": "web"}, true},
		{"missing key", map[string]string{"app": "web"}, map[string]string{"app": "web", "env": "prod"}, false},
		{"wrong value", map[string]string{"app": "api"}, map[string]string{"app": "web"}, false},
		{"empty selector matches all", map[string]string{"app": "web"}, map[string]string{}, true},
		{"nil pod labels", nil, map[string]string{"app": "web"}, false},
		{"nil selector", map[string]string{"app": "web"}, nil, true},
		{"both nil", nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, matchLabels(tt.pod, tt.selector))
		})
	}
}

// ============================================================
// Pod 调度测试
// ============================================================

func TestPodScheduler_SchedulePod(t *testing.T) {
	c := model.NewCluster(10)
	for _, node := range c.Nodes {
		node.MemoryTotal = 8000
		node.ComputePower = 100
	}
	ps := NewPodScheduler(c)

	pod := &Pod{
		ID:       "pod-1",
		Name:     "web-1",
		Phase:    PodPending,
		Resource: model.ResourceRequirement{ComputePower: 10, Memory: 1000},
		Labels:   map[string]string{"app": "web"},
	}

	nodeID, err := ps.SchedulePod(pod)
	require.NoError(t, err)
	assert.NotEmpty(t, nodeID)
	assert.Equal(t, PodRunning, pod.Phase)
}

func TestPodScheduler_NilPod(t *testing.T) {
	ps := NewPodScheduler(model.NewCluster(5))
	_, err := ps.SchedulePod(nil)
	assert.ErrorIs(t, err, ErrNilPod)
}

func TestPodScheduler_AlreadyScheduled(t *testing.T) {
	ps := NewPodScheduler(model.NewCluster(5))
	pod := &Pod{ID: "pod-1", Phase: PodRunning}
	_, err := ps.SchedulePod(pod)
	assert.ErrorIs(t, err, ErrPodAlreadyScheduled)
}

// ============================================================
// ReplicaSet 控制循环测试
// ============================================================

func TestReplicaSet_CreateAndReconcile(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建集群（10 节点）
	// 2. 创建 ReplicaSet（replicas=3）
	// 3. 验证 3 个 Pod 被创建并调度
}

func TestReplicaSet_ScaleUp(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建 ReplicaSet（replicas=2）
	// 2. ScaleReplicaSet → replicas=5
	// 3. Reconcile
	// 4. 验证 Pod 数量从 2 增加到 5
}

func TestReplicaSet_ScaleDown(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建 ReplicaSet（replicas=5）
	// 2. ScaleReplicaSet → replicas=2
	// 3. Reconcile
	// 4. 验证 Pod 数量从 5 减少到 2
}

func TestReplicaSet_SelfHealing(t *testing.T) {
	// TODO(learner): 关键测试！
	// 1. 创建 ReplicaSet（replicas=3）
	// 2. 手动将一个 Pod 标记为 Failed
	// 3. 调用 OnPodFailed
	// 4. Reconcile
	// 5. 验证新 Pod 被创建，总数恢复到 3
}

func TestReplicaSet_NegativeReplicas(t *testing.T) {
	// TODO(learner): 实现
	// ScaleReplicaSet(-1) → ErrInvalidReplicas
}

// ============================================================
// Deployment 滚动更新测试
// ============================================================

func TestDeployment_RollingUpdate(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建 Deployment（replicas=3, v1）
	// 2. 更新 Template（v2）
	// 3. 验证滚动更新过程中：
	//    - 可用 Pod 数 >= replicas - maxUnavailable
	//    - 总 Pod 数 <= replicas + maxSurge
	// 4. 更新完成后所有 Pod 都是 v2
}

func TestDeployment_Rollback(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建 Deployment v1
	// 2. 更新到 v2
	// 3. Rollback
	// 4. 验证所有 Pod 回到 v1 的配置
}

func TestDeployment_RecreateStrategy(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建 Deployment（strategy=Recreate）
	// 2. 更新 Template
	// 3. 验证旧 Pod 全部终止后才创建新 Pod
}

// ============================================================
// Service 服务发现测试
// ============================================================

func TestService_Resolve(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建 Service（selector: app=web）
	// 2. 创建 3 个匹配的 Pod（labels: app=web, phase=Running）
	// 3. ReconcileEndpoints
	// 4. Resolve → 返回某个 Pod ID
	// 5. 多次 Resolve → 应该负载均衡到不同 Pod
}

func TestService_NoEndpoints(t *testing.T) {
	sc := NewServiceController()
	svc := &Service{
		ID: "svc-1", Name: "web", Namespace: "default",
		Selector: map[string]string{"app": "web"},
	}
	sc.CreateService(svc)

	_, err := sc.Resolve("web", "default")
	assert.ErrorIs(t, err, ErrNoEndpoints)
}

func TestService_EndpointReconciliation(t *testing.T) {
	// TODO(learner): 实现
	// 1. 创建 Service
	// 2. 添加匹配的 Pod → Endpoints 更新
	// 3. Pod 变为 Failed → Endpoints 移除该 Pod
	// 4. 验证 Endpoints 只包含 Running 的 Pod
}
