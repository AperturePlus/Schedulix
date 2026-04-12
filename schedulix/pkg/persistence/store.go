package persistence

import (
	"errors"
)

// ─── 错误定义 ───────────────────────────────────────────────

var (
	// ErrKeyNotFound 键不存在。
	ErrKeyNotFound = errors.New("key not found")

	// ErrKeyEmpty 键为空。
	ErrKeyEmpty = errors.New("key must be non-empty")

	// ErrStoreCorrupted 存储数据损坏。
	ErrStoreCorrupted = errors.New("store data is corrupted")

	// ErrStoreClosed 存储已关闭。
	ErrStoreClosed = errors.New("store is closed")

	// ErrWriteFailed 写入失败。
	ErrWriteFailed = errors.New("write operation failed")

	// ErrReadFailed 读取失败。
	ErrReadFailed = errors.New("read operation failed")
)

// ─── 存储接口 ───────────────────────────────────────────────
//
// 学习要点：
//   定义统一的存储接口，底层实现可以是内存、文件、或数据库。
//   这是"依赖倒置原则"的实践：上层模块依赖抽象接口，不依赖具体实现。
//
//   Schedulix 提供三种实现：
//   1. MemoryStore — 内存存储（开发/测试用）
//   2. FileStore — 文件存储（学习文件 I/O）
//   3. WALStore — 预写日志存储（学习崩溃恢复）

// Store 键值存储接口。
// 所有实现必须是线程安全的。
type Store interface {
	// Get 根据键获取值。
	// 键不存在时返回 nil, ErrKeyNotFound。
	Get(key string) ([]byte, error)

	// Put 写入键值对。
	// 键已存在时覆盖（upsert 语义）。
	Put(key string, value []byte) error

	// Delete 删除键值对。
	// 键不存在时静默成功（幂等）。
	Delete(key string) error

	// List 列出所有键（按字典序）。
	List() ([]string, error)

	// Close 关闭存储，释放资源（如文件句柄）。
	// 关闭后的操作返回 ErrStoreClosed。
	// 幂等：多次调用不报错。
	Close() error
}
