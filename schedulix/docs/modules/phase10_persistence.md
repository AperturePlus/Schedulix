# 阶段十：数据持久化（Persistence）

## 学习目标

实现三种持久化方案（内存、文件、WAL），掌握数据持久性、原子写入和崩溃恢复。

## 前置知识

- 阶段一完成（数据模型）
- 了解文件系统基本概念

## 核心概念

### 1. 为什么需要持久化？

内存中的数据，进程一退出就没了。持久化 = 把数据写到磁盘，下次启动还能恢复。

Schedulix 中需要持久化的数据：
- 集群快照（节点状态、拓扑）
- 任务状态（哪些任务在运行、完成、失败）
- 检查点（任务执行进度）
- 事件日志（故障历史）
- 监控指标（历史数据）

### 2. 三种实现的递进关系

```
MemoryStore（最简单）
    ↓ 加上文件 I/O
FileStore（学习文件操作）
    ↓ 加上预写日志
WALStore（学习崩溃恢复）
```

### 3. 原子写入（Atomic Write）

**问题**：直接写文件，写到一半进程崩溃 → 文件损坏（半写状态）。

**解决**：先写临时文件，再原子重命名。

```
步骤 1: 写入 tmp/random.tmp    ← 崩溃？临时文件丢弃，原文件完好
步骤 2: file.Sync()            ← 确保数据落盘
步骤 3: os.Rename(tmp, target)  ← 原子操作！要么完成，要么不做
```

```go
func atomicWrite(path string, data []byte) error {
    tmpPath := path + ".tmp"
    
    // 1. 写临时文件
    f, err := os.Create(tmpPath)
    if err != nil {
        return err
    }
    
    if _, err := f.Write(data); err != nil {
        f.Close()
        os.Remove(tmpPath)  // 清理
        return err
    }
    
    // 2. 同步到磁盘
    if err := f.Sync(); err != nil {
        f.Close()
        os.Remove(tmpPath)
        return err
    }
    f.Close()
    
    // 3. 原子重命名
    return os.Rename(tmpPath, path)
}
```

### 4. 预写日志（Write-Ahead Log, WAL）

**核心思想**：在修改数据之前，先记录"打算做什么"。

```
正常流程：
1. 写 WAL: {op: "put", key: "task-1", value: "...", committed: false}
2. 写数据文件: data/task-1.dat
3. 更新 WAL: {committed: true}

崩溃在步骤 1 后：
→ WAL 有记录但 committed=false
→ 重启时重放：执行 put 操作
→ 数据恢复！

崩溃在步骤 2 后：
→ WAL 有记录，committed=false
→ 重启时重放：再次执行 put（幂等，覆盖写）
→ 数据恢复！

崩溃在步骤 3 后：
→ WAL 已标记 committed=true
→ 重启时跳过
→ 数据完好！
```

### 5. 文件名安全编码

键可能包含 `/`、`\`、`..` 等危险字符。直接用作文件名会导致：
- 路径穿越攻击：`../../etc/passwd`
- 文件系统错误：`key/with/slashes`

解决：对键进行 hex 编码或 URL 编码。

```go
// "task-123" → "7461736b2d313233.dat"
func safeFileName(key string) string {
    return hex.EncodeToString([]byte(key)) + ".dat"
}
```

## 练习任务

1. 打开 `pkg/persistence/memory.go`，实现内存存储（最简单，先热身）
2. 打开 `pkg/persistence/filestore.go`，实现文件存储（学习文件 I/O 和原子写入）
3. 打开 `pkg/persistence/wal.go`，实现 WAL 存储（学习崩溃恢复）
4. 将 CheckpointStore 改为使用 Store 接口（依赖倒置）
5. 将集群快照保存到 FileStore

## 验证

```bash
go test ./pkg/persistence/...

# 测试崩溃恢复（手动）：
# 1. 写入一些数据
# 2. 在 WAL 写入后、数据写入前模拟崩溃（kill 进程）
# 3. 重启，检查数据是否恢复
```
