package persistence

import (
	"sync"
)

// MemoryStore 内存键值存储。
//
// 学习要点：
//   最简单的 Store 实现，用于开发和测试。
//   数据存在内存中，进程退出即丢失。
//
// 鲁棒性设计：
//   - 线程安全（sync.RWMutex）
//   - 关闭后所有操作返回 ErrStoreClosed
//   - Get 返回值的副本（防止调用方修改内部数据）
//   - Put 存储值的副本（防止调用方后续修改影响存储）
type MemoryStore struct {
	data   map[string][]byte
	mu     sync.RWMutex
	closed bool
}

// NewMemoryStore 创建内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string][]byte),
	}
}

// Get 根据键获取值。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. closed → 返回 nil, ErrStoreClosed
// 2. key == "" → 返回 nil, ErrKeyEmpty
// 3. 键不存在 → 返回 nil, ErrKeyNotFound
// 4. 返回值的副本（make + copy），不返回内部引用
// 5. 使用 mu.RLock()
func (m *MemoryStore) Get(key string) ([]byte, error) {
	// TODO: 实现
	panic("not implemented")
}

// Put 写入键值对。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. closed → 返回 ErrStoreClosed
// 2. key == "" → 返回 ErrKeyEmpty
// 3. value == nil → 存储空切片 []byte{}（不存 nil）
// 4. 存储 value 的副本（make + copy）
// 5. 使用 mu.Lock()
func (m *MemoryStore) Put(key string, value []byte) error {
	// TODO: 实现
	panic("not implemented")
}

// Delete 删除键值对。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. closed → 返回 ErrStoreClosed
// 2. key == "" → 返回 ErrKeyEmpty
// 3. 键不存在 → 静默成功（幂等）
func (m *MemoryStore) Delete(key string) error {
	// TODO: 实现
	panic("not implemented")
}

// List 列出所有键（按字典序）。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// 1. closed → 返回 nil, ErrStoreClosed
// 2. 返回排序后的键列表
// 3. 使用 mu.RLock()
func (m *MemoryStore) List() ([]string, error) {
	// TODO: 实现
	panic("not implemented")
}

// Close 关闭存储。
//
// TODO(learner): 实现此方法
// 鲁棒性要求：
// - 幂等：多次调用不报错
// - 清空 data map 释放内存
func (m *MemoryStore) Close() error {
	// TODO: 实现
	panic("not implemented")
}
