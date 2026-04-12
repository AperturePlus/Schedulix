# 需求文档

## 简介

Schedulix 是一个面向学习者的 Serverless 架构项目，使用 Go 语言开发，模拟万卡（GPU）集群调度系统。项目以渐进式学习为核心设计理念，涵盖 Go 语言基础与高级特性、调度算法、事件模拟、容灾与负载均衡、Serverless 架构设计以及容器操作等多个学习维度。所有 GPU 均为软件模拟，无需真实硬件。项目按阶段递进，从基础语法练习逐步过渡到复杂的分布式调度系统。

## 术语表

- **Schedulix**: 本项目的名称，一个模拟万卡调度的 Serverless 学习平台
- **GPU_Node**: 模拟的 GPU 计算节点，包含状态、算力、故障率等属性
- **Cluster**: 由多个 GPU_Node 组成的模拟集群，最大规模可达 10,000 个节点
- **Scheduler**: 调度器，负责将任务分配到可用的 GPU_Node 上执行
- **Task**: 用户提交的计算任务，包含资源需求、优先级、预计执行时间等属性
- **Event_Simulator**: 事件模拟器，用于生成节点故障、网络延迟、资源波动等模拟事件
- **Fault_Event**: 模拟的故障事件，包括节点宕机、网络分区、性能降级等
- **Recovery_Engine**: 容灾恢复引擎，负责在故障发生后执行任务迁移和节点恢复
- **Load_Balancer**: 负载均衡器，确保任务在集群中均匀分布
- **Task_Queue**: 任务队列，管理待调度任务的优先级排序和分发
- **Metrics_Collector**: 指标收集器，收集集群运行状态、调度性能等监控数据
- **Container_Runtime**: 容器运行时模拟层，模拟容器的创建、启动、停止和销毁操作
- **Serverless_Gateway**: Serverless 网关，处理函数触发、路由和自动扩缩容
- **Learning_Module**: 学习模块，项目中按难度递进的功能单元

## 需求

### 需求 1：Go 语言基础数据结构与项目脚手架

**用户故事：** 作为一名 Go 语言学习者，我希望通过构建项目基础数据模型来练习 Go 的基础语法和常见数据结构操作，以便在实际项目中掌握数组、结构体、字典和字符串的使用。

#### 验收标准

1. THE Schedulix SHALL 提供 GPU_Node 结构体定义，包含节点 ID（字符串）、状态（枚举）、算力（整数）、内存容量（整数）和故障率（浮点数）字段
2. THE Schedulix SHALL 提供 Task 结构体定义，包含任务 ID（字符串）、资源需求（结构体）、优先级（整数）、预计执行时间（整数，毫秒）和提交时间（时间戳）字段
3. THE Schedulix SHALL 使用 Go 切片（slice）管理 GPU_Node 集合，支持按状态过滤、按算力排序操作
4. THE Schedulix SHALL 使用 Go 字典（map）实现节点 ID 到 GPU_Node 的快速查找索引
5. THE Schedulix SHALL 提供节点状态的字符串序列化与反序列化功能，支持 JSON 格式
6. FOR ALL 有效的 GPU_Node 对象，序列化后再反序列化 SHALL 产生与原始对象等价的结果（往返属性）
7. WHEN 创建 Cluster 时，THE Schedulix SHALL 根据指定数量初始化 GPU_Node 集合，每个节点具有唯一 ID

### 需求 2：任务队列与优先级管理

**用户故事：** 作为一名学习者，我希望实现一个任务队列系统来学习 Go 的接口（interface）设计和堆（heap）数据结构，以便理解优先级调度的基本原理。

#### 验收标准

1. THE Task_Queue SHALL 实现 Go 标准库 `container/heap` 接口，支持按优先级排序的任务入队和出队操作
2. WHEN 一个 Task 被提交到 Task_Queue 时，THE Task_Queue SHALL 将该 Task 插入到正确的优先级位置
3. WHEN 调用出队操作时，THE Task_Queue SHALL 返回当前优先级最高的 Task
4. THE Task_Queue SHALL 支持查询当前队列长度和队列是否为空的操作
5. WHEN 两个 Task 具有相同优先级时，THE Task_Queue SHALL 按提交时间先后顺序排列（先提交的优先）
6. IF Task_Queue 为空时执行出队操作，THEN THE Task_Queue SHALL 返回明确的错误信息

### 需求 3：基础调度算法

**用户故事：** 作为一名学习者，我希望实现多种基础调度算法，以便理解不同调度策略的优劣和适用场景。

#### 验收标准

1. THE Scheduler SHALL 实现 First-Fit 调度算法，将任务分配到第一个满足资源需求的可用 GPU_Node
2. THE Scheduler SHALL 实现 Best-Fit 调度算法，将任务分配到剩余资源最少但仍满足需求的可用 GPU_Node
3. THE Scheduler SHALL 实现 Round-Robin 调度算法，将任务按轮询方式均匀分配到可用 GPU_Node
4. WHEN 没有可用的 GPU_Node 满足 Task 的资源需求时，THE Scheduler SHALL 将该 Task 保留在 Task_Queue 中并返回资源不足的状态信息
5. THE Scheduler SHALL 通过 Go 接口（interface）定义统一的调度策略抽象，允许不同算法实现可互换
6. WHEN 一个 Task 被成功调度时，THE Scheduler SHALL 更新目标 GPU_Node 的已用资源信息

### 需求 4：Go 并发模型与并行调度

**用户故事：** 作为一名学习者，我希望利用 Go 的 goroutine 和 channel 实现并发调度，以便深入理解 Go 的并发编程模型和常见并发模式。

#### 验收标准

1. THE Scheduler SHALL 使用 goroutine 并发处理多个 Task 的调度请求
2. THE Scheduler SHALL 使用 Go channel 实现 Task 的生产者-消费者模式，生产者提交任务到 channel，消费者从 channel 取出任务进行调度
3. THE Scheduler SHALL 使用 sync.Mutex 或 sync.RWMutex 保护 Cluster 中共享的 GPU_Node 状态数据，防止数据竞争
4. WHEN 多个 goroutine 同时请求调度时，THE Scheduler SHALL 保证每个 GPU_Node 的资源不会被超额分配
5. THE Scheduler SHALL 使用 context.Context 支持调度操作的超时控制和取消功能
6. WHEN 调度操作超时时，THE Scheduler SHALL 释放已预留的资源并返回超时错误
7. THE Schedulix SHALL 使用 sync.WaitGroup 确保所有并发调度任务完成后才进行结果汇总

### 需求 5：事件模拟引擎

**用户故事：** 作为一名学习者，我希望构建一个事件模拟引擎来模拟真实集群中的各种故障和异常情况，以便学习事件驱动编程和概率模型的应用。

#### 验收标准

1. THE Event_Simulator SHALL 支持生成以下类型的 Fault_Event：节点宕机、网络延迟增大、节点性能降级、节点恢复上线
2. THE Event_Simulator SHALL 基于可配置的概率参数生成 Fault_Event，每种事件类型具有独立的发生概率
3. THE Event_Simulator SHALL 支持时间步进模式，每个时间步中根据概率独立判定各节点是否发生故障
4. WHEN 一个节点宕机事件发生时，THE Event_Simulator SHALL 将对应 GPU_Node 的状态更新为不可用，并记录故障发生时间
5. WHEN 一个节点恢复事件发生时，THE Event_Simulator SHALL 将对应 GPU_Node 的状态更新为可用，并记录恢复时间
6. THE Event_Simulator SHALL 提供事件日志功能，按时间顺序记录所有已发生的 Fault_Event
7. THE Event_Simulator SHALL 支持通过 JSON 配置文件设定各类事件的概率参数和模拟时长
8. FOR ALL 有效的事件配置 JSON，解析后再序列化 SHALL 产生与原始配置等价的结果（往返属性）

### 需求 6：容灾恢复与任务迁移

**用户故事：** 作为一名学习者，我希望实现容灾恢复机制，以便理解分布式系统中故障处理和任务迁移的核心概念。

#### 验收标准

1. WHEN 一个 GPU_Node 发生宕机故障时，THE Recovery_Engine SHALL 检测该节点上所有正在执行的 Task
2. WHEN 检测到受影响的 Task 后，THE Recovery_Engine SHALL 将这些 Task 重新提交到 Task_Queue 进行重新调度
3. THE Recovery_Engine SHALL 实现检查点（checkpoint）机制，定期保存 Task 的执行进度
4. WHEN 一个 Task 被迁移到新的 GPU_Node 时，THE Recovery_Engine SHALL 从最近的检查点恢复执行，而非从头开始
5. IF 连续三次迁移同一个 Task 均失败，THEN THE Recovery_Engine SHALL 将该 Task 标记为失败状态并发出告警
6. THE Recovery_Engine SHALL 记录每次故障恢复操作的详细日志，包括故障类型、受影响任务数、恢复耗时

### 需求 7：负载均衡策略

**用户故事：** 作为一名学习者，我希望实现多种负载均衡策略，以便理解大规模集群中资源利用率优化的方法。

#### 验收标准

1. THE Load_Balancer SHALL 实现静态负载均衡策略，基于节点算力权重分配任务
2. THE Load_Balancer SHALL 实现动态负载均衡策略，基于节点实时资源使用率分配任务
3. WHILE Cluster 处于运行状态时，THE Load_Balancer SHALL 周期性地计算各节点的负载指标（CPU 使用率、内存使用率、任务数量）
4. WHEN 某个 GPU_Node 的负载超过可配置的阈值时，THE Load_Balancer SHALL 触发任务迁移，将部分任务转移到低负载节点
5. THE Load_Balancer SHALL 通过 Go 接口定义统一的负载均衡策略抽象，允许不同策略实现可互换
6. THE Metrics_Collector SHALL 记录每个 GPU_Node 的历史负载数据，支持按时间范围查询

### 需求 8：万卡规模集群模拟

**用户故事：** 作为一名学习者，我希望能够模拟 10,000 个 GPU 节点的大规模集群，以便学习大规模系统的性能优化和资源管理技术。

#### 验收标准

1. THE Cluster SHALL 支持创建包含 10,000 个 GPU_Node 的模拟集群
2. THE Cluster SHALL 支持异构节点配置，不同 GPU_Node 具有不同的算力和内存容量
3. THE Cluster SHALL 将节点组织为多层拓扑结构（机架 → 机柜 → 数据中心），模拟真实物理部署
4. WHEN 模拟 10,000 节点集群时，THE Schedulix SHALL 在单机环境下保持调度延迟在可接受范围内（单次调度操作小于 100 毫秒）
5. THE Cluster SHALL 支持通过配置文件定义集群拓扑和节点参数，避免硬编码
6. THE Schedulix SHALL 提供集群状态的快照功能，支持将当前集群状态序列化为 JSON 并从 JSON 恢复
7. FOR ALL 有效的 Cluster 快照，序列化后再反序列化 SHALL 产生与原始集群状态等价的结果（往返属性）

### 需求 9：Serverless 函数框架

**用户故事：** 作为一名学习者，我希望将调度系统的核心功能封装为 Serverless 函数，以便学习 Serverless 架构的设计原则和实现方式。

#### 验收标准

1. THE Serverless_Gateway SHALL 提供 HTTP API 端点，支持提交任务、查询任务状态、查询集群状态等操作
2. THE Serverless_Gateway SHALL 将每个 API 请求路由到对应的处理函数，每个函数独立且无状态
3. WHEN 收到任务提交请求时，THE Serverless_Gateway SHALL 验证请求参数的完整性和合法性
4. IF 请求参数不合法，THEN THE Serverless_Gateway SHALL 返回包含具体错误描述的 HTTP 400 响应
5. THE Serverless_Gateway SHALL 支持函数的冷启动和热启动模拟，记录每次函数调用的启动延迟
6. THE Serverless_Gateway SHALL 实现自动扩缩容逻辑，根据请求并发量动态调整函数实例数量
7. WHILE 没有请求到达时，THE Serverless_Gateway SHALL 将函数实例数量缩减至零（缩容到零）

### 需求 10：容器运行时模拟

**用户故事：** 作为一名学习者，我希望模拟容器的生命周期管理，以便理解容器编排和资源隔离的基本原理。

#### 验收标准

1. THE Container_Runtime SHALL 模拟容器的完整生命周期：创建（Created）→ 启动（Running）→ 停止（Stopped）→ 销毁（Destroyed）
2. THE Container_Runtime SHALL 为每个容器分配模拟的资源配额（CPU 份额、内存限制）
3. WHEN 创建容器时，THE Container_Runtime SHALL 验证宿主 GPU_Node 的剩余资源是否满足容器的资源配额要求
4. IF 宿主 GPU_Node 资源不足，THEN THE Container_Runtime SHALL 拒绝创建容器并返回资源不足错误
5. THE Container_Runtime SHALL 支持在同一个 GPU_Node 上运行多个容器，各容器的资源配额之和不超过节点总资源
6. WHEN 容器状态发生变化时，THE Container_Runtime SHALL 发出状态变更事件，供其他模块订阅

### 需求 11：监控与可视化数据输出

**用户故事：** 作为一名学习者，我希望收集和输出系统运行的关键指标数据，以便学习监控系统设计并直观地观察调度效果。

#### 验收标准

1. THE Metrics_Collector SHALL 收集以下集群级指标：总任务数、已完成任务数、失败任务数、平均调度延迟、集群资源利用率
2. THE Metrics_Collector SHALL 收集以下节点级指标：每个 GPU_Node 的当前负载、已分配任务数、故障次数、累计运行时间
3. THE Metrics_Collector SHALL 以可配置的时间间隔周期性地采集指标数据
4. THE Metrics_Collector SHALL 将指标数据输出为结构化的 JSON 格式，便于外部工具消费
5. WHEN 指标数据被输出时，THE Metrics_Collector SHALL 包含采集时间戳和数据版本号
6. THE Metrics_Collector SHALL 支持将历史指标数据导出为 CSV 格式，便于学习者使用电子表格工具分析

### 需求 12：渐进式学习路径与项目组织

**用户故事：** 作为一名学习者，我希望项目按照由浅入深的学习路径组织，以便我能循序渐进地掌握各项技术。

#### 验收标准

1. THE Schedulix SHALL 将功能组织为以下渐进式 Learning_Module：阶段一（Go 基础与数据模型）、阶段二（任务队列与基础调度）、阶段三（并发编程）、阶段四（事件模拟与容灾）、阶段五（负载均衡与万卡规模）、阶段六（Serverless 架构）、阶段七（容器操作）、阶段八（监控与集成）
2. THE Schedulix SHALL 为每个 Learning_Module 提供独立的 Go 包（package），模块间通过明确定义的接口通信
3. THE Schedulix SHALL 为每个 Learning_Module 提供 README 文档，说明学习目标、前置知识、核心概念和练习任务
4. WHEN 学习者完成一个 Learning_Module 时，THE Schedulix SHALL 提供该模块的单元测试套件，用于验证实现的正确性
5. THE Schedulix SHALL 确保每个 Learning_Module 可以独立编译和测试，不强制依赖后续阶段的模块
6. THE Schedulix SHALL 在项目根目录提供总体学习路线图文档，描述各阶段的依赖关系和建议学习顺序
