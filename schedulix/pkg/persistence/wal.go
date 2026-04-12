package persistence

import (
	"fmt"
	"sync"
	"time"
)

// ─── 预写日志（Write-Ahead Log, WAL）─────────────────────────
//
// 学习要点：
//   WAL 是数据库和分布式系统中保证数据持久性的核心技术。
//   核心思想：在修改数据之前，先把"打算做什么"写入日志文件。
//   如果进程崩溃，重启后可以从日志中恢复未完成的操作。
//
//   流程：
//   1. 收到写请求
//   2. 将操作写入 WAL 文件（追加写，顺序 I/O，非常快）
//   3. 执行实际的数据修改
//   4. 修改成功后，标记 WAL 条目为已完成
//
//   崩溃恢复：
//   1. 读取 WAL 文件
//   2. 找到所有未标记为"已完成"的条目
//   3. 重放这些操作
//
//   这就是为什么叫"预写"日志 — 先写日志，再写数据。

// WALEntry 预写日志条目。
type WALEntry struct {
	Sequence  int64     `json:"seq"`       // 递增序列号
	Timestamp time.Time `json:"ts"`
	Operation string    `json:"op"`        // "put" 或 "delete"
	Key       string    `json:"key"`
	Value     []byte    `json:"value,omitempty"`
	Committed bool      `json:"committed"` // 操作是否已完成
}

// WALStore 基于预写日志的持久化存储。
// 在 FileStore 基础上增加 WAL，保证崩溃一致性。
//
// 鲁棒性设计：
//   - WAL 文件追加写入（顺序 I/O）
//   - 每次写入后 Sync 确保落盘
//   - 启动时自动回放未完成的 WAL 条目
//   - WAL 文件定期压缩（删除已完成的条目）
//   - WAL 文件损坏时跳过损坏条目，恢复尽可能多的数据
type WALStore struct {
	store    *FileStore // 底层文件存储
	walPath  string     // WAL 文件路径
	sequence int64      // 当前序列号
	mu       sync.Mutex
	closed   bool
}

// NewWALStore 创建 WAL 存储。
//
// TODO(learner): 实现此方法
// 步骤：
// 1. 创建底层 FileStore
// 2. 打开或创建 WAL 文件
// 3. 回放未完成的 WAL 条目（崩溃恢复）
// 4. 设置 sequence 为最大已有序列号 + 1
//
// 鲁棒性要求：
// - WAL 文件不存在 → 创建新文件（首次启动）
// - WAL 文件损坏 → 跳过损坏条目，记录警告，继续恢复
// - 底层 FileStore 创建失败 → 包装错误返回
func NewWALStore(baseDir string) (*WALStore, error) {
	// TODO: 实现
	panic("not implemented")
}

// Get 读取值（直接委托给底层 FileStore）。
func (w *WALStore) Get(key string) ([]byte, error) {
	if w.closed {
		return nil, ErrStoreClosed
	}
	return w.store.Get(key)
}

// Put 写入值（先写 WAL，再写数据）。
//
// TODO(learner): 实现此方法
// 步骤（WAL 协议）：
// 1. closed → ErrStoreClosed
// 2. 创建 WALEntry（Committed = false）
// 3. 追加写入 WAL 文件 + Sync
//    - 写入失败 → 返回错误（数据未修改，安全）
// 4. 调用 store.Put(key, value) 执行实际写入
//    - 写入失败 → WAL 中有记录，下次启动会重放（安全）
// 5. 标记 WALEntry 为 Committed
// 6. 递增 sequence
func (w *WALStore) Put(key string, value []byte) error {
	// TODO: 实现
	panic("not implemented")
}

// Delete 删除值（先写 WAL，再删数据）。
//
// TODO(learner): 实现此方法（同 Put 的 WAL 协议）
func (w *WALStore) Delete(key string) error {
	// TODO: 实现
	panic("not implemented")
}

// List 列出所有键。
func (w *WALStore) List() ([]string, error) {
	if w.closed {
		return nil, ErrStoreClosed
	}
	return w.store.List()
}

// Close 关闭 WAL 存储。
//
// TODO(learner): 实现此方法
// 步骤：
// 1. 压缩 WAL（删除已完成的条目）
// 2. 关闭 WAL 文件
// 3. 关闭底层 FileStore
func (w *WALStore) Close() error {
	// TODO: 实现
	panic("not implemented")
}

// recover 回放未完成的 WAL 条目。
//
// TODO(learner): 实现此方法
// 步骤：
// 1. 逐行读取 WAL 文件
// 2. 解析每行为 WALEntry
//    - 解析失败 → 跳过该行，记录警告
// 3. 找到 Committed == false 的条目
// 4. 按 Sequence 顺序重放：
//    - "put" → store.Put(key, value)
//    - "delete" → store.Delete(key)
// 5. 重放成功后标记为 Committed
func (w *WALStore) recover() error {
	// TODO: 实现
	panic("not implemented")
}

// compact 压缩 WAL 文件，删除已完成的条目。
//
// TODO(learner): 实现此方法
// 步骤：
// 1. 读取所有 WAL 条目
// 2. 过滤掉 Committed == true 的条目
// 3. 将剩余条目写入新的 WAL 文件
// 4. 原子替换旧文件
func (w *WALStore) compact() error {
	// TODO: 实现
	panic("not implemented")
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
