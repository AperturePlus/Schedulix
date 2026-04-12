# 阶段一：Go 基础与数据模型

## 学习目标

通过构建 Schedulix 的核心数据模型，掌握 Go 语言的基础语法和常见数据结构操作。

## 前置知识

- 任意一门编程语言的基础经验
- 了解 JSON 格式

## 核心概念

### 1. 结构体（Struct）

Go 的结构体类似于其他语言的 class，但没有继承。

```go
type GPU_Node struct {
    ID           string  // 字段名大写 = 导出（public）
    computeCache int     // 字段名小写 = 未导出（private）
}
```

### 2. 方法（Method）

Go 通过"接收者"将函数绑定到类型：

```go
// 值接收者：不修改原对象
func (n GPU_Node) AvailableMemory() int {
    return n.MemoryTotal - n.MemoryUsed
}

// 指针接收者：可以修改原对象
func (n *GPU_Node) AllocateMemory(amount int) {
    n.MemoryUsed += amount
}
```

**何时用指针接收者？**
- 需要修改接收者的字段
- 结构体较大，避免拷贝开销

### 3. iota 枚举

Go 没有 enum 关键字，用 `const + iota` 模拟：

```go
type NodeStatus int

const (
    NodeStatusIdle     NodeStatus = iota // 0
    NodeStatusBusy                       // 1（自动递增）
    NodeStatusOffline                    // 2
    NodeStatusDegraded                   // 3
)
```

### 4. 切片操作

```go
// 过滤：返回所有 Idle 节点
func filterIdle(nodes []*GPU_Node) []*GPU_Node {
    result := make([]*GPU_Node, 0)
    for _, n := range nodes {
        if n.Status == NodeStatusIdle {
            result = append(result, n)
        }
    }
    return result
}

// 排序：按算力降序
sort.Slice(nodes, func(i, j int) bool {
    return nodes[i].ComputePower > nodes[j].ComputePower
})
```

### 5. Map（字典）

```go
// 创建
nodes := make(map[string]*GPU_Node)

// 写入
nodes["node-0001"] = &GPU_Node{ID: "node-0001"}

// 读取（注意检查是否存在）
node, ok := nodes["node-0001"]
if !ok {
    // 不存在
}

// 删除
delete(nodes, "node-0001")
```

### 6. JSON 自定义序列化

当默认的 JSON 序列化不满足需求时（如枚举值要序列化为字符串），实现 `json.Marshaler` 和 `json.Unmarshaler` 接口：

```go
func (s NodeStatus) MarshalJSON() ([]byte, error) {
    name, ok := nodeStatusNames[s]
    if !ok {
        return nil, fmt.Errorf("unknown status: %d", s)
    }
    return json.Marshal(name) // 将字符串序列化为 JSON
}

func (s *NodeStatus) UnmarshalJSON(data []byte) error {
    var name string
    if err := json.Unmarshal(data, &name); err != nil {
        return err
    }
    val, ok := nodeStatusValues[name]
    if !ok {
        return fmt.Errorf("unknown status: %s", name)
    }
    *s = val
    return nil
}
```

## 练习任务

1. 打开 `pkg/model/status.go`，实现 `NodeStatus` 和 `TaskStatus` 的 JSON 序列化方法
2. 打开 `pkg/model/node.go`，实现 `AvailableMemory()` 和 `CanAccept()` 方法
3. 打开 `pkg/model/cluster.go`，实现 `NewCluster()`、`GetNode()`、`FilterByStatus()`、`SortByComputePower()`
4. 实现 `SnapshotToJSON()` 和 `RestoreFromJSON()`

## 验证

```bash
go test ./pkg/model/...
```
