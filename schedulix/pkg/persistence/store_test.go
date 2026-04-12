package persistence

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Store 接口一致性测试
// 所有 Store 实现都应通过这些测试
// ============================================================

// runStoreTests 对任意 Store 实现运行标准测试套件。
// 这是一种"接口测试"模式：写一次测试，验证所有实现。
func runStoreTests(t *testing.T, newStore func(t *testing.T) Store) {
	t.Run("PutAndGet", func(t *testing.T) {
		s := newStore(t)
		defer s.Close()

		err := s.Put("key1", []byte("value1"))
		require.NoError(t, err)

		val, err := s.Get("key1")
		require.NoError(t, err)
		assert.Equal(t, []byte("value1"), val)
	})

	t.Run("GetNonExistent", func(t *testing.T) {
		s := newStore(t)
		defer s.Close()

		_, err := s.Get("missing")
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})

	t.Run("PutOverwrite", func(t *testing.T) {
		s := newStore(t)
		defer s.Close()

		s.Put("key1", []byte("v1"))
		s.Put("key1", []byte("v2"))

		val, _ := s.Get("key1")
		assert.Equal(t, []byte("v2"), val)
	})

	t.Run("Delete", func(t *testing.T) {
		s := newStore(t)
		defer s.Close()

		s.Put("key1", []byte("v1"))
		err := s.Delete("key1")
		require.NoError(t, err)

		_, err = s.Get("key1")
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})

	t.Run("DeleteIdempotent", func(t *testing.T) {
		s := newStore(t)
		defer s.Close()

		err := s.Delete("nonexistent")
		assert.NoError(t, err) // 幂等
	})

	t.Run("EmptyKey", func(t *testing.T) {
		s := newStore(t)
		defer s.Close()

		err := s.Put("", []byte("v"))
		assert.ErrorIs(t, err, ErrKeyEmpty)

		_, err = s.Get("")
		assert.ErrorIs(t, err, ErrKeyEmpty)
	})

	t.Run("NilValue", func(t *testing.T) {
		s := newStore(t)
		defer s.Close()

		err := s.Put("key1", nil)
		require.NoError(t, err)

		val, err := s.Get("key1")
		require.NoError(t, err)
		assert.NotNil(t, val) // 应存储空切片而非 nil
	})

	t.Run("List", func(t *testing.T) {
		s := newStore(t)
		defer s.Close()

		s.Put("c", []byte("3"))
		s.Put("a", []byte("1"))
		s.Put("b", []byte("2"))

		keys, err := s.List()
		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, keys) // 字典序
	})

	t.Run("OperationsAfterClose", func(t *testing.T) {
		s := newStore(t)
		s.Close()

		_, err := s.Get("key1")
		assert.ErrorIs(t, err, ErrStoreClosed)

		err = s.Put("key1", []byte("v"))
		assert.ErrorIs(t, err, ErrStoreClosed)
	})

	t.Run("CloseIdempotent", func(t *testing.T) {
		s := newStore(t)
		assert.NoError(t, s.Close())
		assert.NoError(t, s.Close()) // 多次关闭不报错
	})

	t.Run("GetReturnsCopy", func(t *testing.T) {
		// TODO(learner): 实现
		// 1. Put("key1", "hello")
		// 2. val, _ := Get("key1")
		// 3. 修改 val[0] = 'X'
		// 4. 再次 Get("key1") → 应该仍然是 "hello"（不受修改影响）
	})
}

// ============================================================
// MemoryStore 测试
// ============================================================

func TestMemoryStore(t *testing.T) {
	runStoreTests(t, func(t *testing.T) Store {
		return NewMemoryStore()
	})
}

// ============================================================
// FileStore 测试
// ============================================================

func TestFileStore(t *testing.T) {
	runStoreTests(t, func(t *testing.T) Store {
		dir := t.TempDir() // 自动创建临时目录，测试结束自动清理
		store, err := NewFileStore(dir)
		require.NoError(t, err)
		return store
	})
}

// TODO(learner): 添加 FileStore 特有的测试
// - TestFileStore_AtomicWrite: 验证写入中途失败不会损坏已有数据
// - TestFileStore_SafeFileName: 验证特殊字符键被正确编码
// - TestFileStore_LargeValue: 验证大数据写入和读取

// ============================================================
// WALStore 测试
// ============================================================

func TestWALStore(t *testing.T) {
	runStoreTests(t, func(t *testing.T) Store {
		dir := t.TempDir()
		store, err := NewWALStore(dir)
		require.NoError(t, err)
		return store
	})
}

// TODO(learner): 添加 WALStore 特有的测试
// - TestWALStore_CrashRecovery: 模拟崩溃后恢复
// - TestWALStore_Compact: 验证压缩后数据完整
