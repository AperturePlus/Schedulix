# 阶段八：监控与集成

## 学习目标

实现指标采集和数据导出，掌握监控系统设计和环形缓冲区数据结构。

## 前置知识

- 阶段一完成（数据模型）
- 了解 CSV 格式

## 核心概念

### 1. 环形缓冲区（Ring Buffer）

固定大小的数组，写指针循环移动。满时覆盖最旧数据。

```
初始状态（容量 4）：
[_, _, _, _]  cursor=0

写入 A, B, C, D：
[A, B, C, D]  cursor=0（回到开头）

写入 E（覆盖 A）：
[E, B, C, D]  cursor=1

写入 F（覆盖 B）：
[E, F, C, D]  cursor=2
```

**为什么用环形缓冲区？**
- 固定内存占用，不会因为长时间运行而 OOM
- 自动丢弃最旧数据，保留最近的 N 条记录
- 写入 O(1)，读取 O(1)

```go
type RingBuffer struct {
    data   []ClusterMetrics
    size   int
    cursor int
    count  int // 已写入总数（可能 > size）
}

func (rb *RingBuffer) Write(m ClusterMetrics) {
    rb.data[rb.cursor] = m
    rb.cursor = (rb.cursor + 1) % rb.size
    rb.count++
}

func (rb *RingBuffer) GetLatest(n int) []ClusterMetrics {
    if n > rb.size { n = rb.size }
    if n > rb.count { n = rb.count }
    
    result := make([]ClusterMetrics, n)
    for i := 0; i < n; i++ {
        idx := (rb.cursor - n + i + rb.size) % rb.size
        result[i] = rb.data[idx]
    }
    return result
}
```

### 2. 指标采集

周期性地遍历集群状态，计算汇总指标：

```go
func (mc *MetricsCollector) Collect() {
    cluster := mc.cluster
    
    // 集群级指标
    metrics := ClusterMetrics{
        Timestamp: time.Now(),
        Version:   mc.version,
    }
    mc.version++
    
    totalMem := 0
    usedMem := 0
    for _, node := range cluster.Nodes {
        totalMem += node.MemoryTotal
        usedMem += node.MemoryUsed
        // ... 累加其他指标
    }
    
    metrics.ResourceUtilization = float64(usedMem) / float64(totalMem)
    
    // 写入环形缓冲区
    mc.buffer.Write(metrics)
}
```

### 3. CSV 流式导出

使用 `encoding/csv` 逐行写入，避免将全部数据加载到内存：

```go
import "encoding/csv"

func exportCSV(w io.Writer, data []ClusterMetrics) error {
    writer := csv.NewWriter(w)
    defer writer.Flush()
    
    // 写表头
    writer.Write([]string{
        "timestamp", "version", "total_tasks", 
        "completed_tasks", "failed_tasks",
        "avg_schedule_delay_ms", "resource_utilization",
    })
    
    // 逐行写数据
    for _, m := range data {
        writer.Write([]string{
            m.Timestamp.Format(time.RFC3339),
            strconv.Itoa(m.Version),
            strconv.Itoa(m.TotalTasks),
            // ...
        })
    }
    
    return writer.Error()
}
```

## 练习任务

1. 打开 `pkg/metrics/collector.go`，实现环形缓冲区和指标采集
2. 打开 `pkg/metrics/exporter.go`，实现 JSON 和 CSV 导出

## 验证

```bash
go test ./pkg/metrics/...
```
