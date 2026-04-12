# Implementation Plan: Schedulix

## Overview

Schedulix 是一个面向学习者的 Go 语言 Serverless GPU 集群调度模拟器。实现按照 8 个渐进式学习阶段组织，从 Go 基础数据模型逐步构建到完整的万卡规模调度系统。每个阶段对应独立的 Go 包，通过接口解耦，支持独立编译和测试。

## Tasks

- [ ] 1. 项目初始化与数据模型层（阶段一：pkg/model）
  - [ ] 1.1 初始化 Go 模块与项目目录结构
    - 创建 `go.mod`（module 名称 `schedulix`，Go 1.21+）
    - 创建目录结构：`cmd/server/`、`pkg/model/`、`pkg/queue/`、`pkg/scheduler/`、`pkg/simulator/`、`pkg/recovery/`、`pkg/balancer/`、`pkg/gateway/`、`pkg/container/`、`pkg/metrics/`、`configs/`、`docs/modules/`
    - 添加依赖：`pgregory.net/rapid`、`github.com/stretchr/testify`
    - _Requirements: 12.1, 12.2, 12.5_

  - [ ] 1.2 实现 GPU_Node 结构体与状态枚举
    - 在 `pkg/model/node.go` 中定义 `NodeStatus`（int + iota：Idle, Busy, Offline, Degraded）
    - 实现 `NodeStatus` 的 `MarshalJSON`/`UnmarshalJSON` 方法（字符串序列化）
    - 定义 `GPU_Node` 结构体（ID, Status, ComputePower, MemoryTotal, MemoryUsed, FaultRate, AssignedTasks, FaultCount, UptimeMs, RackID, CabinetID, DataCenterID）
    - 实现 `AvailableMemory()` 和 `CanAccept(req)` 方法
    - _Requirements: 1.1, 1.5_

  - [ ] 1.3 实现 Task 结构体与资源需求模型
    - 在 `pkg/model/task.go` 中定义 `TaskStatus`（int + iota：Pending, Running, Completed, Failed, Migrating）及其 JSON 序列化
    - 定义 `ResourceRequirement` 结构体（ComputePower, Memory）
    - 定义 `Task` 结构体（ID, Resource, Priority, EstimatedTimeMs, SubmitTime, Status, AssignedNodeID, MigrationCount, Progress）
    - _Requirements: 1.2_

  - [ ] 1.4 实现 Cluster 结构体与拓扑模型
    - 在 `pkg/model/cluster.go` 中定义 `Cluster`（Nodes map, DataCenters, statusIndex, mu RWMutex）
    - 定义 `DataCenter`、`Cabinet`、`Rack` 拓扑结构体
    - 实现 `NewCluster(count int)` 工厂函数，根据指定数量初始化节点集合，每个节点唯一 ID
    - 实现 `GetAvailableNodes(status)` 使用辅助索引返回节点列表
    - 实现 `UpdateNodeStatus(nodeID, newStatus)` 同步更新辅助索引
    - 实现按状态过滤、按算力排序的切片操作辅助函数
    - _Requirements: 1.3, 1.4, 1.7, 8.3_

  - [ ] 1.5 实现 Checkpoint、Container、FaultEvent 等辅助模型
    - 在 `pkg/model/checkpoint.go` 中定义 `Checkpoint` 结构体
    - 在 `pkg/model/container.go` 中定义 `ContainerState` 枚举和 `Container` 结构体
    - 在 `pkg/simulator/event.go` 中定义 `FaultType` 枚举和 `FaultEvent` 结构体
    - 在 `pkg/simulator/config.go` 中定义 `EventConfig` 结构体及 JSON 解析验证
    - _Requirements: 5.1, 5.7, 6.3, 10.1_

  - [ ] 1.6 实现 JSON 序列化与反序列化功能
    - 实现 `GPU_Node` 的完整 JSON 序列化/反序列化
    - 实现 `Cluster` 的快照序列化/反序列化（`SnapshotToJSON` / `RestoreFromJSON`），恢复后重建 statusIndex 和 mu
    - 实现 `EventConfig` 的 JSON 解析与验证（概率值范围 [0.0, 1.0]）
    - _Requirements: 1.5, 1.6, 5.7, 5.8, 8.6_

  - [ ]* 1.7 属性测试：GPU_Node 序列化往返一致性
    - **Property 1: GPU_Node Round-Trip Consistency**
    - 使用 `pgregory.net/rapid` 生成随机 GPU_Node，验证 JSON 序列化后反序列化产生等价对象
    - **Validates: Requirements 1.6**

  - [ ]* 1.8 属性测试：EventConfig 序列化往返一致性
    - **Property 2: EventConfig Round-Trip Consistency**
    - 使用 `rapid` 生成随机 EventConfig（概率值 [0.0, 1.0]），验证 JSON 往返一致性
    - **Validates: Requirements 5.8**

  - [ ]* 1.9 属性测试：Cluster 快照往返一致性
    - **Property 3: Cluster Snapshot Round-Trip Consistency**
    - 使用 `rapid` 生成随机 Cluster（含拓扑），验证快照序列化后反序列化产生等价集群状态
    - **Validates: Requirements 8.7**

  - [ ]* 1.10 单元测试：数据模型层
    - 测试 NodeStatus JSON 序列化/反序列化（各状态值）
    - 测试 `CanAccept` 方法（正常、资源不足、Offline、Degraded 场景）
    - 测试 `NewCluster` 节点唯一性
    - 测试 `GetAvailableNodes` 和 `UpdateNodeStatus` 索引一致性
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.7_

- [ ] 2. Checkpoint — 阶段一验证
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 3. 任务队列与基础调度（阶段二：pkg/queue, pkg/scheduler）
  - [ ] 3.1 实现优先级任务队列
    - 在 `pkg/queue/priority_queue.go` 中定义内部 `taskHeap` 类型，实现 `container/heap` 的五个接口方法（Len, Less, Swap, Push, Pop）
    - Less 逻辑：优先级高的在前；同优先级按 SubmitTime 先提交的在前
    - 封装线程安全的 `TaskQueue` 结构体（内部 `sync.Mutex`），实现 `Enqueue`、`Dequeue`、`Peek`、`Len`、`IsEmpty` 方法
    - `Dequeue` 和 `Peek` 在队列为空时返回 `ErrQueueEmpty` 错误
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_

  - [ ]* 3.2 单元测试：优先级队列
    - 测试按优先级出队顺序
    - 测试同优先级按提交时间 FIFO
    - 测试空队列出队返回错误
    - 测试 Len 和 IsEmpty
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_

  - [ ] 3.3 定义调度策略接口与实现 First-Fit 算法
    - 在 `pkg/scheduler/strategy.go` 中定义 `ScheduleStrategy` 接口（Schedule, Name 方法）
    - 定义 `ErrNoAvailableNode` 错误
    - 在 `pkg/scheduler/firstfit.go` 中实现 First-Fit：遍历可用节点，返回第一个满足资源需求的节点 ID
    - _Requirements: 3.1, 3.5_

  - [ ] 3.4 实现 Best-Fit 和 Round-Robin 调度算法
    - 在 `pkg/scheduler/bestfit.go` 中实现 Best-Fit：遍历可用节点，选择剩余资源最少但满足需求的节点
    - 在 `pkg/scheduler/roundrobin.go` 中实现 Round-Robin：维护游标索引，轮询分配到可用节点
    - _Requirements: 3.2, 3.3_

  - [ ] 3.5 实现调度器核心逻辑
    - 创建 `Scheduler` 结构体，持有 `ScheduleStrategy` 接口引用、`TaskQueue` 引用和 `Cluster` 引用
    - 实现 `ScheduleTask` 方法：从队列取任务 → 调用策略选择节点 → 更新节点已用资源（MemoryUsed, AssignedTasks）→ 更新任务状态为 Running
    - 无可用节点时将任务保留在队列并返回资源不足状态
    - 支持运行时切换调度策略
    - _Requirements: 3.4, 3.5, 3.6_

  - [ ]* 3.6 单元测试：调度算法
    - 测试 First-Fit 选择第一个满足节点
    - 测试 Best-Fit 选择剩余资源最少的节点
    - 测试 Round-Robin 轮询均匀分配
    - 测试无可用节点时返回错误
    - 测试调度后节点资源更新正确
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.6_

- [ ] 4. Checkpoint — 阶段二验证
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 5. 并发调度（阶段三：pkg/scheduler/concurrent）
  - [ ] 5.1 实现生产者-消费者并发调度模型
    - 在 `pkg/scheduler/concurrent.go` 中实现 `ConcurrentScheduler`
    - 使用 buffered channel（容量 1000）作为任务通道
    - 实现生产者方法：将任务发送到 channel
    - 实现消费者 goroutine：从 channel 接收任务，调用 `ScheduleStrategy` 执行调度
    - 使用 `sync.Mutex` / `sync.RWMutex` 保护 Cluster 共享状态，防止数据竞争
    - _Requirements: 4.1, 4.2, 4.3_

  - [ ] 5.2 实现超时控制与资源安全
    - 使用 `context.Context` 支持调度操作的超时和取消
    - 超时时释放已预留资源并返回超时错误
    - 使用 `sync.WaitGroup` 确保所有并发调度任务完成后汇总结果
    - 保证多 goroutine 并发调度时节点资源不被超额分配
    - _Requirements: 4.4, 4.5, 4.6, 4.7_

  - [ ]* 5.3 单元测试：并发调度
    - 测试多 goroutine 并发调度不产生数据竞争（使用 `go test -race`）
    - 测试超时取消正确释放资源
    - 测试 WaitGroup 等待所有任务完成
    - 测试资源不超额分配
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7_

- [ ] 6. Checkpoint — 阶段三验证
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 7. 事件模拟与容灾恢复（阶段四：pkg/simulator, pkg/recovery）
  - [ ] 7.1 实现事件模拟引擎
    - 在 `pkg/simulator/engine.go` 中实现 `EventSimulator`
    - 实现离散事件队列（按时间排序的优先级队列）
    - 实现时间步进模式：每步对每个节点独立进行伯努利试验，根据 `EventConfig` 概率生成 FaultEvent
    - 支持四种事件类型：节点宕机、网络延迟、性能降级、节点恢复
    - 节点宕机时更新状态为 Offline 并记录故障时间；恢复时更新为 Idle 并记录恢复时间
    - 实现事件日志功能，按时间顺序记录所有 FaultEvent
    - 支持通过 JSON 配置文件设定概率参数和模拟时长
    - 实现 `EventHandler` 接口注册与回调通知机制
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7_

  - [ ]* 7.2 单元测试：事件模拟引擎
    - 测试各类故障事件正确更新节点状态
    - 测试事件日志按时间排序
    - 测试 EventConfig JSON 解析与验证
    - 测试 EventHandler 回调触发
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6_

  - [ ] 7.3 实现容灾恢复引擎
    - 在 `pkg/recovery/engine.go` 中实现 `RecoveryEngine`
    - 实现故障检测：节点宕机时检测该节点上所有正在执行的 Task
    - 实现任务重提交：将受影响 Task 重新提交到 TaskQueue
    - 实现三次迁移失败标记：`MigrationCount >= 3` 时标记任务为 Failed 并发出告警
    - 记录每次恢复操作的详细日志（故障类型、受影响任务数、恢复耗时）
    - _Requirements: 6.1, 6.2, 6.5, 6.6_

  - [ ] 7.4 实现检查点机制
    - 在 `pkg/recovery/checkpoint.go` 中实现检查点管理
    - 使用内存 `map[string]*Checkpoint` 存储检查点
    - 实现定期保存 Task 执行进度（TaskID, Progress, Timestamp, NodeID）
    - 任务迁移时从最近检查点恢复执行进度
    - _Requirements: 6.3, 6.4_

  - [ ]* 7.5 单元测试：容灾恢复
    - 测试故障检测正确识别受影响任务
    - 测试任务重提交到队列
    - 测试三次迁移失败标记
    - 测试检查点保存与恢复
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_

- [ ] 8. Checkpoint — 阶段四验证
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 9. 负载均衡与万卡规模（阶段五：pkg/balancer, pkg/model 扩展）
  - [ ] 9.1 实现负载均衡策略
    - 在 `pkg/balancer/strategy.go` 中定义 `BalanceStrategy` 接口（SelectNode, ShouldRebalance, Name）
    - 在 `pkg/balancer/static.go` 中实现静态负载均衡：基于节点算力权重分配任务
    - 在 `pkg/balancer/dynamic.go` 中实现动态负载均衡：基于实时资源使用率分配任务
    - 实现阈值触发迁移：节点负载超过可配置阈值时，将部分任务转移到低负载节点
    - _Requirements: 7.1, 7.2, 7.4, 7.5_

  - [ ] 9.2 实现周期性负载指标计算
    - 在 `pkg/balancer/dynamic.go` 或 `pkg/metrics/collector.go` 中实现周期性计算各节点负载指标（CPU 使用率、内存使用率、任务数量）
    - 记录每个 GPU_Node 的历史负载数据，支持按时间范围查询
    - _Requirements: 7.3, 7.6_

  - [ ]* 9.3 单元测试：负载均衡
    - 测试静态策略按权重分配
    - 测试动态策略按实时负载分配
    - 测试阈值触发迁移逻辑
    - 测试 ShouldRebalance 判断
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

  - [ ] 9.4 实现万卡规模集群支持
    - 扩展 `NewCluster` 支持创建 10,000 节点集群，支持异构节点配置（不同算力和内存）
    - 实现多层拓扑结构初始化（DataCenter → Cabinet → Rack → Nodes）
    - 支持通过 JSON 配置文件定义集群拓扑和节点参数
    - 实现按状态分组的辅助索引（`map[NodeStatus][]string`），调度时只遍历目标状态节点
    - _Requirements: 8.1, 8.2, 8.3, 8.5_

  - [ ] 9.5 实现集群快照功能
    - 实现集群状态序列化为 JSON（`SnapshotToJSON`）
    - 实现从 JSON 恢复集群状态（`RestoreFromJSON`），重建 statusIndex 和 mu
    - _Requirements: 8.6, 8.7_

  - [ ]* 9.6 单元测试与性能验证：万卡规模
    - 测试 10,000 节点集群创建
    - 测试异构节点配置
    - 测试拓扑结构正确性
    - 测试单次调度操作延迟 < 100ms（benchmark test）
    - 测试集群快照序列化/反序列化
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6_

- [ ] 10. Checkpoint — 阶段五验证
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 11. Serverless 网关（阶段六：pkg/gateway）
  - [ ] 11.1 实现 HTTP 路由与请求处理
    - 在 `pkg/gateway/router.go` 中使用 `http.ServeMux` 注册 API 路由
    - 在 `pkg/gateway/handler.go` 中实现各 API handler：
      - `POST /api/v1/tasks` — 提交任务（验证参数完整性和合法性，不合法返回 HTTP 400 + 错误描述）
      - `GET /api/v1/tasks/{id}` — 查询任务状态
      - `GET /api/v1/cluster/status` — 查询集群状态
      - `GET /api/v1/cluster/nodes` — 查询节点列表
      - `POST /api/v1/cluster/snapshot` — 创建集群快照
      - `PUT /api/v1/cluster/snapshot` — 从快照恢复
      - `POST /api/v1/simulator/start` — 启动事件模拟
      - `GET /api/v1/metrics` — 获取监控指标
      - `GET /api/v1/metrics/export` — 导出历史指标（支持 `?format=json|csv`）
    - 每个 handler 无状态，通过依赖注入访问 Cluster、Scheduler 等组件
    - _Requirements: 9.1, 9.2, 9.3, 9.4_

  - [ ] 11.2 实现冷/热启动模拟与自动扩缩容
    - 在 `pkg/gateway/scaler.go` 中实现冷启动模拟（`time.Sleep` 可配置延迟）和热启动判断（`sync.Once` 或时间戳）
    - 记录每次函数调用的启动延迟
    - 实现自动扩缩容：维护函数实例计数器，根据并发请求数动态调整
    - 无请求时通过 `time.AfterFunc` 延迟缩容到零
    - _Requirements: 9.5, 9.6, 9.7_

  - [ ]* 11.3 单元测试：Serverless 网关
    - 测试各 API 端点的正常响应
    - 测试参数验证与 HTTP 400 错误响应
    - 测试冷启动/热启动延迟记录
    - 测试扩缩容逻辑（实例增减、缩容到零）
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5, 9.6, 9.7_

- [ ] 12. Checkpoint — 阶段六验证
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 13. 容器运行时（阶段七：pkg/container）
  - [ ] 13.1 实现容器生命周期状态机
    - 在 `pkg/container/lifecycle.go` 中实现容器状态机（Created → Running → Stopped → Destroyed）
    - 非法状态转换返回 `ErrInvalidStateTransition` 错误
    - 实现 `ContainerLifecycle` 观察者接口，状态变更时通知所有订阅者
    - _Requirements: 10.1, 10.6_

  - [ ] 13.2 实现容器运行时资源管理
    - 在 `pkg/container/runtime.go` 中实现 `ContainerRuntime`
    - 创建容器时验证宿主 GPU_Node 剩余资源是否满足容器资源配额（CPUShares, MemoryLimit）
    - 资源不足时拒绝创建并返回错误
    - 支持同一 GPU_Node 上运行多个容器，各容器资源配额之和不超过节点总资源
    - _Requirements: 10.2, 10.3, 10.4, 10.5_

  - [ ]* 13.3 单元测试：容器运行时
    - 测试完整生命周期状态转换
    - 测试非法状态转换返回错误
    - 测试资源配额验证（充足/不足）
    - 测试多容器资源累加不超限
    - 测试观察者通知回调
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5, 10.6_

- [ ] 14. Checkpoint — 阶段七验证
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 15. 监控与数据导出（阶段八：pkg/metrics）
  - [ ] 15.1 实现指标采集器
    - 在 `pkg/metrics/collector.go` 中实现 `MetricsCollector`
    - 采集集群级指标：总任务数、已完成任务数、失败任务数、平均调度延迟、集群资源利用率
    - 采集节点级指标：每个 GPU_Node 的当前负载、已分配任务数、故障次数、累计运行时间
    - 以可配置时间间隔周期性采集
    - 使用环形缓冲区（ring buffer）存储最近 N 个指标快照，避免内存无限增长
    - 每次采集包含时间戳和递增版本号
    - _Requirements: 11.1, 11.2, 11.3, 11.5_

  - [ ] 15.2 实现指标导出器
    - 在 `pkg/metrics/exporter.go` 中实现 `MetricsExporter`
    - 实现 JSON 格式导出（结构化输出）
    - 实现 CSV 格式导出（流式写入，使用 `encoding/csv`）
    - _Requirements: 11.4, 11.6_

  - [ ]* 15.3 单元测试：监控系统
    - 测试指标采集正确性（集群级和节点级）
    - 测试环形缓冲区覆盖最旧数据
    - 测试版本号递增
    - 测试 JSON 和 CSV 导出格式
    - _Requirements: 11.1, 11.2, 11.3, 11.4, 11.5, 11.6_

- [ ] 16. Checkpoint — 阶段八验证
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 17. 集成组装与学习路径文档
  - [ ] 17.1 实现 Serverless 网关入口与依赖注入
    - 在 `cmd/server/main.go` 中组装所有组件：创建 Cluster → 初始化 TaskQueue → 创建 Scheduler → 初始化 EventSimulator、RecoveryEngine、LoadBalancer、ContainerRuntime、MetricsCollector → 注册到 Gateway → 启动 HTTP 服务
    - 通过依赖注入将各组件连接，确保无循环依赖
    - _Requirements: 9.1, 12.1, 12.2_

  - [ ] 17.2 创建配置文件示例
    - 创建 `configs/cluster.json`：集群拓扑配置示例（含 DataCenter、Cabinet、Rack、异构节点参数）
    - 创建 `configs/events.json`：事件模拟配置示例（各类故障概率、模拟步数、步间隔）
    - _Requirements: 5.7, 8.5_

  - [ ] 17.3 创建学习路径文档
    - 创建 `docs/roadmap.md`：总体学习路线图，描述 8 个阶段的依赖关系和建议学习顺序
    - 为每个阶段创建 `docs/modules/phase{1-8}.md` README，说明学习目标、前置知识、核心概念和练习任务
    - _Requirements: 12.1, 12.3, 12.6_

  - [ ]* 17.4 集成测试：端到端流程验证
    - 测试完整任务提交 → 调度 → 执行 → 完成流程
    - 测试故障注入 → 容灾恢复 → 任务迁移流程
    - 测试监控指标采集与导出
    - _Requirements: 12.4, 12.5_

- [ ] 18. Final Checkpoint — 全部验证
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation at each learning phase boundary
- Property tests (1.7, 1.8, 1.9) validate universal round-trip correctness properties
- Unit tests validate specific examples and edge cases
- All code uses Go 1.21+ with standard library + `pgregory.net/rapid` + `github.com/stretchr/testify`
