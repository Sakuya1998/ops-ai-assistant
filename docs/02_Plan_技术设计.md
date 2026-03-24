# 基于Eino框架的运维AI助手 - Plan（技术设计方案）

> 版本: v1.1  
> 更新日期: 2026-03-24  
> 状态: 待评审  
> 对齐文档: Spec v1.1

---

## 1. 设计目标

本方案用于支撑 Spec v1.1 中的 P0/P1 功能落地，确保系统具备以下能力：
- 告警驱动的自动分诊与证据化根因分析
- 对话驱动的跨源诊断
- 可控自动化修复（审批、回滚、审计）
- RAG知识复用与复盘沉淀

关键设计约束：
- 以 Go + Eino 为核心编排框架
- 高危动作必须满足“强制审批 + 全量审计”
- 所有结论均需可追溯证据来源
- 面向 Kubernetes 环境优先优化

---

## 2. 需求到能力映射

| Spec条目 | 技术能力 | 对应模块 |
|---------|---------|---------|
| US-1 自然语言诊断 | 意图识别、实体抽取、结构化响应 | Conversation Orchestrator |
| US-2 告警根因定位 | 告警聚合、指标日志关联、假设排序 | Alert Workflow + RCA Engine |
| US-3 自动修复 | Playbook执行、风险门禁、回滚 | Action Engine + Policy Guard |
| US-4 知识库问答 | 检索、重排、引用回传 | RAG Pipeline |
| US-5 多智能体协作 | Host-Specialist并行诊断 | Multi-Agent Controller |

---

## 3. 总体架构

### 3.1 四层架构

```text
┌────────────────────────────────────────────────────────────────────┐
│                           Interface Layer                          │
│  Web UI / CLI / Alert Webhook / OpenAPI                           │
├────────────────────────────────────────────────────────────────────┤
│                        Orchestration Layer                         │
│  Eino Graph Runtime + Router + Multi-Agent Controller             │
├────────────────────────────────────────────────────────────────────┤
│                          Capability Layer                          │
│  ChatModel / Tool Runtime / Retriever / Policy Guard / Audit      │
├────────────────────────────────────────────────────────────────────┤
│                            Data Layer                              │
│  Prometheus / Loki / Kubernetes API / VikingDB / Incident Store   │
└────────────────────────────────────────────────────────────────────┘
```

### 3.2 逻辑组件

| 组件 | 职责 | 输入 | 输出 |
|------|------|------|------|
| API Gateway | 鉴权、限流、路由 | HTTP/Webhook请求 | 统一内部请求 |
| Conversation Orchestrator | 会话编排与上下文管理 | 用户问题/上下文 | 诊断任务 |
| Alert Workflow | 告警解析、聚合、分派 | 告警事件 | Incident草稿 |
| RCA Engine | 证据收集、根因评分 | 指标/日志/事件 | Top3根因假设 |
| Action Engine | 执行修复动作与验证 | 批准动作 | 执行结果 |
| Knowledge Pipeline | 文档切片、向量检索、引用 | 问题/检索请求 | 引用化答案 |
| Audit Service | 审计落库与追踪 | 执行事件 | 审计记录 |

---

## 4. 核心流程设计

### 4.1 告警驱动分诊流程（P0）

```text
Alert Webhook
  -> Alert Parser
  -> Dedup/Aggregator
  -> Metrics + Logs + Events Collect
  -> RCA Engine (Top3 Hypothesis)
  -> Decision Node (Auto / Confirm / Escalate)
  -> Incident Report + Notification
```

关键点：
- 聚合键：`service + alert_rule + env + 10min_window`
- 分诊窗口：默认最近 30 分钟
- 输出格式：状态级别、影响面、证据链、建议动作

### 4.2 对话驱动诊断流程（P0）

```text
User Query
  -> Intent & Entity Parsing
  -> Tool Plan Generation
  -> Parallel Tool Calls
  -> Evidence Fusion
  -> Structured Response
```

关键点：
- 实体抽取至少包含：服务名、环境、时间窗、指标类型
- 对话保留上下文窗口：最近 N 轮摘要 + 关键证据缓存

### 4.3 修复执行与回滚流程（P1）

```text
Action Request
  -> Policy Guard (RBAC + Risk)
  -> Approval Gate (for High Risk)
  -> Pre-check Snapshot
  -> Execute Playbook
  -> Post-check Validation
  -> Success or Rollback
  -> Audit Record + Notification
```

关键点：
- 高风险动作必须审批，不支持绕过
- 失败自动回滚，回滚同样写审计

---

## 5. Eino编排设计

### 5.1 Graph节点定义

| 节点名 | 类型 | 功能 |
|-------|------|------|
| `parse_input` | Node | 解析告警或用户意图 |
| `route_domain` | Branch | 选择诊断域（网络/DB/应用/基础设施） |
| `collect_metrics` | Node | 查询 Prometheus 指标 |
| `collect_logs` | Node | 查询 Loki 日志 |
| `collect_events` | Node | 查询 K8s 事件 |
| `rca_rank` | Node | 根因假设生成与排序 |
| `risk_gate` | Branch | 判断自动执行/审批/升级 |
| `execute_action` | Node | 执行Playbook |
| `validate_result` | Node | 修复后健康验证 |
| `publish_report` | Node | 输出报告与通知 |

### 5.2 关键编排伪代码

```go
graph := compose.NewGraph[InputMessage, OpsResult]()
graph.AddNode("parse_input", parseInputNode)
graph.AddBranch("route_domain", routeDomainBranch)
graph.AddNode("collect_metrics", collectMetricsNode)
graph.AddNode("collect_logs", collectLogsNode)
graph.AddNode("collect_events", collectEventsNode)
graph.AddNode("rca_rank", rcaRankNode)
graph.AddBranch("risk_gate", riskGateBranch)
graph.AddNode("execute_action", executeActionNode)
graph.AddNode("validate_result", validateResultNode)
graph.AddNode("publish_report", publishReportNode)
```

---

## 6. 多智能体设计（Host-Specialist）

### 6.1 Agent角色

| Agent | 关注领域 | 主要工具 |
|------|----------|----------|
| SupervisorAgent | 任务拆解、协同、汇总 | route、merge、policy |
| NetworkAgent | 连通性、延迟、丢包 | ping、trace、k8s svc/event |
| DBAgent | 连接池、慢查询、锁等待 | db metrics、db logs |
| AppAgent | 应用错误、依赖超时 | app logs、trace |
| InfraAgent | 节点资源、调度异常 | node/pod metrics、k8s events |

### 6.2 协作协议

- Supervisor 负责分配 `TaskID` 与 `EvidenceScope`
- Specialist 必须返回结构化结果：`conclusion + confidence + evidence`
- Aggregator 对冲突结论执行“证据权重优先”策略
- 最终报告输出统一可信度分层

---

## 7. 工具层设计

### 7.1 工具清单与风险分级

| 工具 | 说明 | 风险级别 | 审批策略 |
|------|------|---------|---------|
| `query_metrics` | Prometheus 即时/区间查询 | 低 | 无需审批 |
| `query_logs` | Loki 检索日志 | 低 | 无需审批 |
| `get_k8s_resource` | 获取资源描述 | 低 | 无需审批 |
| `get_k8s_events` | 获取事件信息 | 低 | 无需审批 |
| `restart_pod` | 重启单个Pod | 中 | ops_engineer及以上 |
| `scale_deployment` | 调整副本数 | 高 | ops_admin审批 |
| `exec_runbook` | 执行修复剧本 | 中/高 | 按剧本风险判定 |
| `rollback_action` | 执行回滚 | 中/高 | 跟随原动作策略 |

### 7.2 工具调用约束

- 每次调用必须携带 `request_id` 与 `incident_id`
- 默认超时：读操作 10s，写操作 30s
- 失败重试：最多 2 次，指数退避
- 写操作必须先经过 `Policy Guard`

---

## 8. 数据模型设计

### 8.1 核心实体

```go
type Incident struct {
    ID            string
    Title         string
    Severity      string
    Status        string
    SourceAlertID string
    Service       string
    Environment   string
    StartedAt     time.Time
    ResolvedAt    *time.Time
}

type Hypothesis struct {
    ID          string
    IncidentID  string
    Domain      string
    Description string
    Confidence  float64
}

type Evidence struct {
    ID          string
    IncidentID  string
    Type        string
    Source      string
    Summary     string
    CapturedAt  time.Time
}

type ActionRecord struct {
    ID           string
    IncidentID   string
    ActionType   string
    RiskLevel    string
    ApprovedBy   *string
    Status       string
    StartedAt    time.Time
    FinishedAt   *time.Time
}
```

### 8.2 存储策略

| 数据类型 | 存储介质 | 保留策略 |
|---------|---------|---------|
| Incident元数据 | PostgreSQL | 1年 |
| 审计日志 | PostgreSQL + 对象存储归档 | 3年 |
| 日志原文引用 | Loki | 依赖日志平台策略 |
| 指标快照摘要 | 对象存储 | 90天 |
| 知识向量 | VikingDB | 持续保留，按版本重建 |

---

## 9. API与事件契约

### 9.1 外部API

| 接口 | 方法 | 用途 |
|------|------|------|
| `/api/v1/chat` | POST | 对话式诊断 |
| `/api/v1/alerts/webhook` | POST | 告警接入 |
| `/api/v1/incidents/{id}` | GET | 获取事件详情 |
| `/api/v1/actions/{id}/approve` | POST | 审批高危动作 |
| `/api/v1/actions/{id}/execute` | POST | 执行动作 |
| `/api/v1/reports/{id}` | GET | 获取诊断报告 |

### 9.2 Webhook契约（简化）

```json
{
  "alert_id": "string",
  "rule_name": "string",
  "severity": "critical|high|medium|low",
  "service": "payment-service",
  "environment": "prod",
  "starts_at": "2026-03-24T10:00:00Z",
  "labels": {},
  "annotations": {}
}
```

---

## 10. 安全与治理

### 10.1 权限模型

| 角色 | 能力边界 |
|------|---------|
| `ops_reader` | 只读查询、查看报告 |
| `ops_engineer` | 中风险动作执行 |
| `ops_admin` | 高风险动作审批与执行 |

### 10.2 安全控制点

- 网关鉴权：OIDC/JWT + 细粒度RBAC
- 风险门禁：动作风险分级 + 审批流
- 输出脱敏：密钥、令牌、密码字段遮蔽
- 网络边界：仅白名单访问 K8s API 与内部数据源
- 审计合规：操作全链路留痕与不可篡改归档

---

## 11. 可观测性与SLO落地

### 11.1 监控指标

| 指标名 | 说明 | 目标 |
|-------|------|------|
| `ops_ai_chat_latency_p95` | 对话响应P95 | < 5s |
| `ops_ai_triage_latency_p95` | 告警初判P95 | < 30s |
| `ops_ai_rca_top3_hit_rate` | 根因Top3命中率 | ≥ 75% |
| `ops_ai_action_success_rate` | 修复动作成功率 | ≥ 50% |
| `ops_ai_highrisk_misexec_count` | 高危误执行次数 | 0 |

### 11.2 追踪与日志

- 全链路 TraceID：请求、工具调用、审批、动作执行统一串联
- 结构化日志字段：`trace_id`, `incident_id`, `action_id`, `agent`
- 异常分级：平台异常、数据源异常、模型异常、策略拒绝

---

## 12. 部署与容量规划

### 12.1 部署拓扑

```text
Kubernetes Namespace: ops-ai
  - ops-ai-api (3 replicas)
  - ops-ai-worker (2 replicas)
  - ops-ai-web (2 replicas)
  - redis (cache/queue)
  - postgres (metadata/audit)
  - vikingdb (vector retrieval)
```

### 12.2 建议资源配额

| 组件 | 副本 | CPU | 内存 |
|------|------|-----|------|
| ops-ai-api | 3 | 2 vCPU | 4 Gi |
| ops-ai-worker | 2 | 4 vCPU | 8 Gi |
| ops-ai-web | 2 | 0.5 vCPU | 1 Gi |
| redis | 3 | 1 vCPU | 2 Gi |
| postgres | 2 | 2 vCPU | 4 Gi |
| vikingdb | 3 | 2 vCPU | 8 Gi |

### 12.3 高可用与灾备

- API/Worker 多副本 + Pod 反亲和
- Postgres 主备或云托管高可用
- 核心数据每日备份，RPO ≤ 15min，RTO ≤ 30min

---

## 13. 测试与发布策略

### 13.1 测试分层

| 层级 | 目标 | 样例 |
|------|------|------|
| 单元测试 | 节点逻辑正确性 | 实体解析、风险判定 |
| 集成测试 | 工具调用与编排正确性 | 告警到报告链路 |
| 回放测试 | 真实历史事件回放 | MTTR与命中率评估 |
| 混沌演练 | 异常场景韧性 | 数据源超时、模型降级 |

### 13.2 发布策略

- Alpha：内部小流量，验证链路稳定性
- Beta：生产影子告警，默认人工审批
- GA：核心业务线启用，按周评估SLO

---

## 14. 里程碑与交付

| 里程碑 | 交付内容 | 验收口径 |
|-------|---------|---------|
| M1 | 编排骨架与对话诊断基础链路 | 能完成结构化诊断回答 |
| M2 | 告警接入 + 自动分诊 | Webhook告警可自动输出Top3 |
| M3 | 多智能体并行诊断 | 至少3类Agent并行可用 |
| M4 | 修复执行与审批回滚 | 高危动作审批闭环可用 |
| M5 | RAG知识库上线 | 回答带引用来源 |
| M6 | 生产就绪与治理 | 审计、监控、SLO看板齐备 |

---

## 15. 风险与应对

| 风险 | 影响 | 应对策略 |
|------|------|---------|
| 模型幻觉导致误导诊断 | 高 | 强制证据化输出 + 低置信度提醒 |
| 外部数据源超时 | 中 | 并行超时控制 + 降级策略 |
| 工具权限过大 | 高 | 最小权限 + 审批门禁 + 审计 |
| 告警洪峰压垮系统 | 高 | 去重聚合 + 限流 + 优先级队列 |
| 检索质量波动 | 中 | 多路检索 + 重排 + 反馈闭环 |

---

## 16. 变更记录

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|---------|------|
| v1.1 | 2026-03-24 | 对齐 Spec v1.1，补充接口契约、治理、SLO、测试发布方案 | AI Assistant |
| v1.0 | 2026-02-13 | 初始版本 | AI Assistant |

---

*文档结束*
