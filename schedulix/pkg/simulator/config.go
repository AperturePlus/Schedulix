package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ─── 错误定义 ───────────────────────────────────────────────

var (
	// ErrInvalidConfig 配置不合法。
	ErrInvalidConfig = errors.New("invalid event config")

	// ErrConfigParseFailed JSON 解析失败。
	ErrConfigParseFailed = errors.New("failed to parse event config JSON")
)

// EventConfig 事件模拟配置（JSON 可序列化）。
//
// 概率说明：
//   - 每个概率值取值范围 [0.0, 1.0]
//   - 各事件类型独立判定，概率之和不要求为 1.0
//   - 例如 NodeDownProb=0.01 表示每步每节点有 1% 概率宕机
type EventConfig struct {
	NodeDownProb     float64 `json:"node_down_prob"`     // 每步每节点宕机概率
	NetworkDelayProb float64 `json:"network_delay_prob"` // 每步每节点网络延迟概率
	DegradedProb     float64 `json:"degraded_prob"`      // 每步每节点性能降级概率
	RecoveryProb     float64 `json:"recovery_prob"`      // 每步每离线节点恢复概率
	TotalSteps       int     `json:"total_steps"`        // 模拟总步数
	StepIntervalMs   int64   `json:"step_interval_ms"`   // 每步间隔（毫秒）
}

// DefaultEventConfig 返回安全的默认配置。
// 当用户配置不合法时，可以 fallback 到此默认值。
func DefaultEventConfig() *EventConfig {
	return &EventConfig{
		NodeDownProb:     0.005,
		NetworkDelayProb: 0.01,
		DegradedProb:     0.008,
		RecoveryProb:     0.05,
		TotalSteps:       100,
		StepIntervalMs:   100,
	}
}

// Validate 验证配置参数的合法性。
//
// TODO(learner): 实现此方法
// 鲁棒性要求 — 检查所有字段，收集所有错误（不要遇到第一个就返回）：
// 1. 所有概率值在 [0.0, 1.0] 范围内
//    - 超出范围 → 收集错误信息 "node_down_prob (1.5) out of range [0.0, 1.0]"
// 2. TotalSteps > 0
// 3. StepIntervalMs > 0
// 4. 如果有多个错误，用 fmt.Errorf 拼接所有错误信息返回
//    — 这样调用方一次就能看到所有问题，而不是修一个报一个
//
// 提示：可以用 []string 收集错误消息，最后 join
func (c *EventConfig) Validate() error {
	// TODO: 实现
	panic("not implemented")
}

// ParseConfig 从 JSON 字节解析事件配置。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. data 为 nil 或空 → 返回 DefaultEventConfig() 和一个警告错误
//    （优雅降级：给你默认值，但告诉你输入有问题）
// 2. 使用 json.Unmarshal — 失败 → 包装为 ErrConfigParseFailed
// 3. 调用 Validate() — 失败 → 包装为 ErrInvalidConfig
// 4. 全部通过 → 返回配置和 nil
func ParseConfig(data []byte) (*EventConfig, error) {
	// TODO: 实现
	panic("not implemented")
}

// Clamp 将配置中的概率值钳制到 [0.0, 1.0] 范围。
// 这是一种"自动修复"策略：与其拒绝不合法的值，不如修正它。
//
// TODO(learner): 实现此方法
// 提示：math.Max(0.0, math.Min(1.0, value))
func (c *EventConfig) Clamp() {
	// TODO: 实现
	panic("not implemented")
}

// --- 防止 unused import ---
var _ = json.Unmarshal
var _ = fmt.Sprintf
