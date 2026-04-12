package persistence

import (
	"fmt"
	"sync"
)

// FileStore 基于文件系统的键值存储。
//
// 学习要点：
//   每个键对应一个文件，文件名 = 键的 hash 或安全编码。
//   学习 Go 的文件 I/O：os.Create, os.ReadFile, os.WriteFile, os.Remove
//
//   存储布局：
//     baseDir/
//     ├── metadata.json    # 存储元数据（版本、键列表）
//     ├── data/
//     │   ├── <key1>.dat   # 键值数据文件
//     │   ├── <key2>.dat
//     │   └── ...
//     └── tmp/             # 临时文件目录（原子写入用）
//
// 鲁棒性设计：
//   - 原子写入：先写临时文件，再 rename（防止写入中途崩溃导致数据损坏）
//   - 文件名安全编码：键中的特殊字符（/、\、..）被编码，防止路径穿越攻击
//   - 关闭后所有操作返回 ErrStoreClosed
//   - 目录不存在时自动创建
//   - 读取损坏文件时返回 ErrStoreCorrupted 而非 panic
type FileStore struct {
	baseDir string
	mu      sync.RWMutex
	closed  bool
}

// NewFileStore 创建文件存储。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. baseDir == "" → 返回 nil, 错误
// 2. 创建 baseDir、baseDir/data、baseDir/tmp 目录（os.MkdirAll）
//    - 创建失败 → 包装错误返回
// 3. 加载或创建 metadata.json
func NewFileStore(baseDir string) (*FileStore, error) {
	// TODO: 实现
	panic("not implemented")
}

// Get 从文件读取值。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. closed → ErrStoreClosed
// 2. key == "" → ErrKeyEmpty
// 3. 安全编码 key → 文件路径
// 4. 文件不存在 → ErrKeyNotFound
// 5. 读取失败 → 包装为 ErrReadFailed
// 6. 使用 mu.RLock()
func (fs *FileStore) Get(key string) ([]byte, error) {
	// TODO: 实现
	panic("not implemented")
}

// Put 原子写入值到文件。
//
// TODO(learner): 实现此方法
// 鲁棒性要求（原子写入三步法）：
// 1. closed → ErrStoreClosed
// 2. key == "" → ErrKeyEmpty
// 3. 写入临时文件：baseDir/tmp/<random>.tmp
//    - 写入失败 → 清理临时文件，返回 ErrWriteFailed
// 4. 同步到磁盘：file.Sync()（确保数据落盘）
// 5. 原子重命名：os.Rename(tmpPath, dataPath)
//    - rename 失败 → 清理临时文件，返回 ErrWriteFailed
// 6. 使用 mu.Lock()
//
// 为什么要原子写入？
//   如果直接写目标文件，写到一半进程崩溃 → 文件损坏。
//   先写临时文件再 rename，rename 是原子操作 → 要么完整写入，要么不写。
func (fs *FileStore) Put(key string, value []byte) error {
	// TODO: 实现
	panic("not implemented")
}

// Delete 删除键对应的文件。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. closed → ErrStoreClosed
// 2. key == "" → ErrKeyEmpty
// 3. 文件不存在 → 静默成功（幂等）
// 4. 删除失败 → 包装错误返回
func (fs *FileStore) Delete(key string) error {
	// TODO: 实现
	panic("not implemented")
}

// List 列出所有键。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. 读取 baseDir/data 目录下的所有 .dat 文件
// 2. 将文件名反编码为键
// 3. 排序后返回
// 4. 目录读取失败 → 包装错误返回
func (fs *FileStore) List() ([]string, error) {
	// TODO: 实现
	panic("not implemented")
}

// Close 关闭文件存储。
//
// TODO(learner): 实现此方法
func (fs *FileStore) Close() error {
	// TODO: 实现
	panic("not implemented")
}

// safeFileName 将键编码为安全的文件名。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - 替换 /、\、..、空格等特殊字符
// - 可以用 hex 编码或 URL 编码
// - 确保编码后的文件名不超过 255 字符（文件系统限制）
func safeFileName(key string) string {
	// TODO: 实现
	panic("not implemented")
}

// unsafeFileName 将安全文件名解码回原始键。
//
// TODO(learner): 实现此方法
func unsafeFileName(filename string) string {
	// TODO: 实现
	panic("not implemented")
}

// --- 防止 unused import ---
var _ = fmt.Sprintf
