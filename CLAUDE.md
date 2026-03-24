# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 项目概览

这是一个基于 Eino（字节跳动开源的 Go 语言 AI 框架）构建的运维 AI 助手文档库。系统为运维团队提供智能告警分诊、根因分析、自动化修复和知识库问答能力。

**技术栈：**
- AI 框架：Eino (Go)
- LLM：DeepSeek-V3
- 向量数据库：VikingDB
- 可观测性：Prometheus + Loki
- 部署：Kubernetes

## 架构设计

### 四层架构

1. **接口层**：Web UI、CLI、告警 Webhook、OpenAPI
2. **编排层**：Eino Graph 运行时、路由器、多智能体控制器
3. **能力层**：ChatModel、工具运行时、检索器、策略门禁、审计
4. **数据层**：Prometheus、Loki、Kubernetes API、VikingDB、事件存储

### 多智能体系统（Host-Specialist 模式）

- **SupervisorAgent**：任务拆解、协调、汇总
- **NetworkAgent**：连通性、延迟、丢包诊断
- **DBAgent**：连接池、慢查询、锁分析
- **AppAgent**：应用错误、依赖超时
- **InfraAgent**：节点资源、调度异常

## 核心工作流

### 告警驱动分诊（P0）
告警 Webhook → 解析器 → 去重/聚合器 → 指标/日志/事件收集 → RCA 引擎（Top3 假设）→ 决策节点 → 事件报告

### 对话驱动诊断（P0）
用户查询 → 意图与实体解析 → 工具计划 → 并行工具调用 → 证据融合 → 结构化响应

### 修复与回滚（P1）
动作请求 → 策略门禁（RBAC + 风险）→ 审批门禁 → 预检快照 → 执行 Playbook → 后检验证 → 成功/回滚 → 审计记录

## Eino Graph 节点

编排图中的关键节点：
- `parse_input`：解析告警或用户意图
- `route_domain`：选择诊断域（网络/数据库/应用/基础设施）
- `collect_metrics`、`collect_logs`、`collect_events`：并行数据收集
- `rca_rank`：根因假设生成与排序
- `risk_gate`：判断自动执行/审批/升级
- `execute_action`：执行 Playbook
- `validate_result`：修复后健康检查
- `publish_report`：输出报告与通知

## 工具风险分级

- **低风险**（无需审批）：`query_metrics`、`query_logs`、`get_k8s_resource`、`get_k8s_events`
- **中风险**（ops_engineer 及以上）：`restart_pod`、`exec_runbook`（中风险）
- **高风险**（ops_admin 审批）：`scale_deployment`、`exec_runbook`（高风险）

所有写操作必须通过策略门禁并生成审计记录。

## 数据模型

核心实体：`Incident`、`Hypothesis`、`Evidence`、`ActionRecord`

存储策略：
- 事件元数据：PostgreSQL（保留 1 年）
- 审计日志：PostgreSQL + 对象存储归档（保留 3 年）
- 日志引用：Loki（依赖日志平台策略）
- 指标快照：对象存储（保留 90 天）
- 知识向量：VikingDB（持续保留）

## 关键 SLO

- 对话响应 P95：< 5s
- 告警分诊 P95：< 30s
- RCA Top3 命中率：≥ 75%
- 动作成功率：≥ 50%
- 高危误执行次数：0

## 安全与治理

- 所有执行类工具必须进行 RBAC 验证
- 高风险动作必须审批 + 审计追踪
- 输出必须脱敏（不含密钥、令牌、凭据）
- 全链路可追溯：TraceID 串联请求、工具调用、审批、动作

## 文档结构

本仓库包含 11 个结构化文档，覆盖完整生命周期：

1. **01_Spec_产品需求.md**：产品需求、用户故事、验收标准
2. **02_Plan_技术设计.md**：技术架构、Eino 编排、多智能体设计
3. **03_API_接口契约.md**：API 契约、认证、错误码
4. **04_Deployment_部署手册.md**：部署流程、配置、回滚
5. **05_Runbook_运行手册.md**：值班 SOP、故障处置、升级路径
6. **06_Test_测试策略.md**：测试策略、关键用例、质量门槛
7. **07_Security_安全合规.md**：权限模型、审批门禁、审计要求
8. **08_SRE_SLO与告警.md**：SLO 指标、告警分级、值班规则
9. **09_Agent_Prompt策略.md**：多智能体角色、Prompt 治理
10. **10_Eval_评测标准.md**：评测指标、数据集、通过标准
11. **11_KB_知识库治理.md**：知识来源、入库规范、检索质量

**推荐阅读路径：**
- 产品/项目经理：01 → 02 → 10
- 架构师/研发工程师：02 → 03 → 04 → 06
- SRE/运维值班：05 → 07 → 08 → 11
