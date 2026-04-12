package gateway

import (
	"sync"
	"sync/atomic"
	"time"
)

// Scaler Serverless 自动扩缩容模拟器。
//
// 鲁棒性设计：
//   - 使用 atomic 操作保证计数器线程安全
//   - activeInstances 不会变为负数（clamp to 0）
//   - 冷启动延迟为 0 或负值时跳过 Sleep
//   - 缩容定时器取消逻辑防止竞态
type Scaler struct {
	activeInstances atomic.Int64
	coldStartDelay  time.Duration
	scaleDownDelay  time.Duration
	lastRequestTime atomic.Int64 // Unix 毫秒时间戳
	scaleDownTimer  *time.Timer
	mu              sync.Mutex
}

// NewScaler 创建扩缩容模拟器。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - coldStartDelay < 0 → 设为 0（不信任调用方）
// - scaleDownDelay < 0 → 设为 0
func NewScaler(coldStartDelay, scaleDownDelay time.Duration) *Scaler {
	// TODO: 实现
	panic("not implemented")
}

// OnRequest 请求到达时调用，处理扩容和冷启动。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 取消正在进行的缩容定时器（新请求到达，不应缩容）
// 2. 如果 activeInstances == 0（冷启动）：
//    - coldStartDelay > 0 时才 Sleep
// 3. activeInstances.Add(1)
// 4. 更新 lastRequestTime
// 5. 返回是否为冷启动
func (s *Scaler) OnRequest() (isColdStart bool) {
	// TODO: 实现
	panic("not implemented")
}

// OnRequestDone 请求完成时调用，处理缩容。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. activeInstances.Add(-1)
// 2. 如果结果 < 0 → 修正为 0（防御性：不应发生，但以防万一）
// 3. 如果 activeInstances == 0 且 scaleDownDelay > 0：
//    - 启动 time.AfterFunc 延迟缩容
//    - 缩容回调中再次检查 activeInstances == 0（可能在延迟期间有新请求）
func (s *Scaler) OnRequestDone() {
	// TODO: 实现
	panic("not implemented")
}

// GetActiveInstances 返回当前活跃实例数。
func (s *Scaler) GetActiveInstances() int64 {
	return s.activeInstances.Load()
}
