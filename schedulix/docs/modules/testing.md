# Go 测试完全指南

## 为什么测试是核心技能？

代码写完不算完，测试通过才算完。测试是你和未来的自己（以及队友）之间的契约：
- **验证**（Verification）：代码做了它该做的事吗？
- **回归保护**（Regression）：改了一处，别处没坏吧？
- **文档**（Documentation）：测试就是最好的使用示例
- **设计反馈**（Design Feedback）：难以测试的代码 = 设计有问题

## Go 测试基础

### 文件命名约定

```
node.go          ← 源代码
node_test.go     ← 测试代码（必须以 _test.go 结尾）
```

测试文件和源代码在同一个目录、同一个包中。

### 测试函数签名

```go
import "testing"

// 函数名必须以 Test 开头，参数必须是 *testing.T
func TestAvailableMemory(t *testing.T) {
    node := &GPU_Node{MemoryTotal: 1000, MemoryUsed: 300}
    got := node.AvailableMemory()
    want := 700
    if got != want {
        t.Errorf("AvailableMemory() = %d, want %d", got, want)
    }
}
```

### 运行测试

```bash
# 测试单个包
go test ./pkg/model/...

# 测试所有包
go test ./...

# 显示详细输出
go test -v ./pkg/model/...

# 运行特定测试
go test -run TestAvailableMemory ./pkg/model/...

# 带竞争检测（并发测试必用）
go test -race ./pkg/scheduler/...

# 性能基准测试
go test -bench=. ./pkg/model/...

# 测试覆盖率
go test -cover ./pkg/model/...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 六种测试模式

### 1. 基础单元测试

最简单的测试：一个输入，一个期望输出。

```go
func TestCanAccept_Idle(t *testing.T) {
    node := &GPU_Node{
        Status:       NodeStatusIdle,
        ComputePower: 100,
        MemoryTotal:  8000,
        MemoryUsed:   2000,
    }
    req := ResourceRequirement{ComputePower: 50, Memory: 4000}
    
    if !node.CanAccept(req) {
        t.Error("expected idle node with sufficient resources to accept task")
    }
}
```

### 2. 表驱动测试（Table-Driven Tests）

Go 社区最推崇的模式。一张表覆盖所有场景：

```go
func TestCanAccept(t *testing.T) {
    tests := []struct {
        name   string
        node   *GPU_Node
        req    ResourceRequirement
        want   bool
    }{
        {
            name: "idle node with enough resources",
            node: &GPU_Node{Status: NodeStatusIdle, ComputePower: 100, MemoryTotal: 8000, MemoryUsed: 2000},
            req:  ResourceRequirement{ComputePower: 50, Memory: 4000},
            want: true,
        },
        {
            name: "offline node always rejects",
            node: &GPU_Node{Status: NodeStatusOffline, ComputePower: 100, MemoryTotal: 8000},
            req:  ResourceRequirement{ComputePower: 1, Memory: 1},
            want: false,
        },
        {
            name: "degraded node halves compute power",
            node: &GPU_Node{Status: NodeStatusDegraded, ComputePower: 100, MemoryTotal: 8000},
            req:  ResourceRequirement{ComputePower: 60, Memory: 1000},
            want: false, // 100/2=50 < 60
        },
        {
            name: "insufficient memory",
            node: &GPU_Node{Status: NodeStatusIdle, ComputePower: 100, MemoryTotal: 8000, MemoryUsed: 7500},
            req:  ResourceRequirement{ComputePower: 10, Memory: 1000},
            want: false, // 8000-7500=500 < 1000
        },
        {
            name: "negative resource request",
            node: &GPU_Node{Status: NodeStatusIdle, ComputePower: 100, MemoryTotal: 8000},
            req:  ResourceRequirement{ComputePower: -1, Memory: 1000},
            want: false, // 不信任输入
        },
        {
            name: "nil node",
            node: nil,
            req:  ResourceRequirement{ComputePower: 10, Memory: 1000},
            want: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 对 nil node 的特殊处理
            var got bool
            if tt.node != nil {
                got = tt.node.CanAccept(tt.req)
            }
            if got != tt.want {
                t.Errorf("CanAccept() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**为什么表驱动？**
- 新增测试用例只需加一行，不需要写新函数
- 所有用例格式统一，一目了然
- `t.Run` 让每个用例有独立的名字，失败时精确定位

### 3. 错误路径测试

测试不只是验证"正确的输入产生正确的输出"，更要验证"错误的输入产生正确的错误"：

```go
func TestDequeue_EmptyQueue(t *testing.T) {
    q := NewTaskQueue()
    
    task, err := q.Dequeue()
    
    if task != nil {
        t.Error("expected nil task from empty queue")
    }
    if !errors.Is(err, ErrQueueEmpty) {
        t.Errorf("expected ErrQueueEmpty, got %v", err)
    }
}

func TestEnqueue_NilTask(t *testing.T) {
    q := NewTaskQueue()
    
    err := q.Enqueue(nil)
    
    if !errors.Is(err, ErrNilTask) {
        t.Errorf("expected ErrNilTask, got %v", err)
    }
}
```

### 4. 属性测试（Property-Based Testing）

不是测试具体的输入输出，而是测试"对于所有合法输入，某个属性都成立"。

使用 `pgregory.net/rapid` 库：

```go
import "pgregory.net/rapid"

// 属性：序列化后反序列化，得到的对象与原始对象等价
func TestGPUNode_JSONRoundTrip(t *testing.T) {
    rapid.Check(t, func(t *rapid.T) {
        // rapid 自动生成随机 GPU_Node
        node := &GPU_Node{
            ID:           rapid.String().Draw(t, "id"),
            Status:       NodeStatus(rapid.IntRange(0, 3).Draw(t, "status")),
            ComputePower: rapid.IntRange(1, 1000).Draw(t, "power"),
            MemoryTotal:  rapid.IntRange(1, 100000).Draw(t, "memTotal"),
            MemoryUsed:   rapid.IntRange(0, 100000).Draw(t, "memUsed"),
            FaultRate:    rapid.Float64Range(0, 1).Draw(t, "faultRate"),
        }
        
        // 序列化
        data, err := json.Marshal(node)
        if err != nil {
            t.Fatal(err)
        }
        
        // 反序列化
        var restored GPU_Node
        err = json.Unmarshal(data, &restored)
        if err != nil {
            t.Fatal(err)
        }
        
        // 验证属性：往返一致
        if node.ID != restored.ID {
            t.Errorf("ID mismatch: %s != %s", node.ID, restored.ID)
        }
        if node.ComputePower != restored.ComputePower {
            t.Errorf("ComputePower mismatch")
        }
        // ... 检查所有字段
    })
}
```

**属性测试 vs 单元测试：**
- 单元测试：你选几个例子来测
- 属性测试：框架自动生成成百上千个随机输入来测
- 属性测试能发现你想不到的边界情况

### 5. 并发测试 + 竞争检测

```go
func TestTaskQueue_ConcurrentEnqueueDequeue(t *testing.T) {
    q := NewTaskQueue()
    const numTasks = 1000
    
    // 并发入队
    var wg sync.WaitGroup
    for i := 0; i < numTasks; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            task := &model.Task{
                ID:       fmt.Sprintf("task-%d", id),
                Priority: id % 10,
            }
            if err := q.Enqueue(task); err != nil {
                t.Errorf("Enqueue failed: %v", err)
            }
        }(i)
    }
    wg.Wait()
    
    // 验证数量
    if q.Len() != numTasks {
        t.Errorf("queue length = %d, want %d", q.Len(), numTasks)
    }
    
    // 并发出队
    results := make(chan *model.Task, numTasks)
    for i := 0; i < numTasks; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            task, err := q.Dequeue()
            if err != nil {
                return // 队列可能被其他 goroutine 清空
            }
            results <- task
        }()
    }
    wg.Wait()
    close(results)
    
    // 验证没有重复
    seen := make(map[string]bool)
    for task := range results {
        if seen[task.ID] {
            t.Errorf("duplicate task: %s", task.ID)
        }
        seen[task.ID] = true
    }
}
```

**必须用 `-race` 运行：**
```bash
go test -race ./pkg/queue/...
```

### 6. 基准测试（Benchmark）

测量性能，验证万卡规模下的延迟要求：

```go
func BenchmarkFirstFit_10000Nodes(b *testing.B) {
    cluster := model.NewCluster(10000)
    strategy := &FirstFitStrategy{}
    task := &model.Task{
        Resource: model.ResourceRequirement{ComputePower: 10, Memory: 1024},
    }
    
    b.ResetTimer() // 不计入初始化时间
    for i := 0; i < b.N; i++ {
        strategy.Schedule(task, cluster)
    }
}
// 输出：BenchmarkFirstFit_10000Nodes-8  50000  25000 ns/op
// 意思：每次调度约 25μs（远小于 100ms 要求）
```

## 测试辅助技巧

### testify 断言库

比 `if got != want` 更简洁：

```go
import "github.com/stretchr/testify/assert"

func TestSomething(t *testing.T) {
    assert.Equal(t, 700, node.AvailableMemory())
    assert.True(t, node.CanAccept(req))
    assert.Nil(t, err)
    assert.ErrorIs(t, err, ErrQueueEmpty)
    assert.NotNil(t, cluster)
    assert.Len(t, nodes, 10)
}
```

### t.Helper()

标记辅助函数，让错误报告指向调用方而非辅助函数内部：

```go
func assertNodeAccepts(t *testing.T, node *GPU_Node, req ResourceRequirement) {
    t.Helper() // 关键！
    if !node.CanAccept(req) {
        t.Errorf("expected node %s to accept request", node.ID)
    }
}
```

### t.Cleanup()

测试结束后自动清理资源：

```go
func TestFileStore(t *testing.T) {
    dir := t.TempDir() // 自动创建临时目录，测试结束自动删除
    store, err := NewFileStore(dir)
    assert.NoError(t, err)
    t.Cleanup(func() {
        store.Close()
    })
    // ... 测试逻辑
}
```

### t.Parallel()

标记测试可以并行运行（加速测试套件）：

```go
func TestA(t *testing.T) {
    t.Parallel()
    // ...
}

func TestB(t *testing.T) {
    t.Parallel()
    // ...
}
```

## 测试金字塔

```
        /  E2E 测试  \        ← 少量，慢，验证完整流程
       / 集成测试      \       ← 中量，验证组件间交互
      / 单元测试          \    ← 大量，快，验证单个函数
```

Schedulix 中：
- **单元测试**：每个方法的正常/异常路径（`*_test.go`）
- **集成测试**：Scheduler + Queue + Cluster 协作（`scheduler_integration_test.go`）
- **E2E 测试**：HTTP API → 调度 → 故障 → 恢复 完整流程
