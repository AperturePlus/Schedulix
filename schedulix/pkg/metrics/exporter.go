package metrics

import "io"

// MetricsExporter 指标导出器。
// 支持 JSON 和 CSV 两种格式。
type MetricsExporter struct {
	collector *MetricsCollector
}

// NewMetricsExporter 创建导出器。
func NewMetricsExporter(collector *MetricsCollector) *MetricsExporter {
	return &MetricsExporter{collector: collector}
}

// ExportJSON 将历史指标导出为 JSON 格式。
//
// TODO(learner): 实现此方法
// 提示：
// 1. 从 collector 获取历史数据
// 2. 使用 encoding/json.NewEncoder(w).Encode() 写入
func (e *MetricsExporter) ExportJSON(w io.Writer) error {
	// TODO: 实现
	panic("not implemented")
}

// ExportCSV 将历史指标导出为 CSV 格式。
//
// TODO(learner): 实现此方法
// 提示：
// 1. 使用 encoding/csv.NewWriter(w)
// 2. 先写表头行
// 3. 逐行写入数据
// 4. 调用 writer.Flush()
func (e *MetricsExporter) ExportCSV(w io.Writer) error {
	// TODO: 实现
	panic("not implemented")
}
