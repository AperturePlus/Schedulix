package gateway

import "net/http"

// Router Serverless 网关路由器。
// 使用标准库 http.ServeMux 注册 API 路由。
type Router struct {
	mux     *http.ServeMux
	handler *Handler
	scaler  *Scaler
}

// NewRouter 创建路由器并注册所有 API 端点。
//
// TODO(learner): 实现此方法
// 注册以下路由：
//   POST /api/v1/tasks          → handler.SubmitTask
//   GET  /api/v1/tasks/         → handler.GetTaskStatus (从 URL 解析 task ID)
//   GET  /api/v1/cluster/status → handler.GetClusterStatus
//   GET  /api/v1/cluster/nodes  → handler.GetNodes
//   POST /api/v1/simulator/start → handler.StartSimulator
//   GET  /api/v1/metrics        → handler.GetMetrics
//   GET  /api/v1/metrics/export → handler.ExportMetrics
//
// 提示：使用 mux.HandleFunc(pattern, handlerFunc)
func NewRouter(handler *Handler, scaler *Scaler) *Router {
	// TODO: 实现
	panic("not implemented")
}

// ServeHTTP 实现 http.Handler 接口。
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
