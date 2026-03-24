# 运维AI助手

基于Eino框架的运维AI助手需求与方案文档集合，覆盖从产品定义到技术落地的完整链路。

---

## 技术栈

- **AI框架**: Eino（字节跳动开源Go框架）
- **LLM**: DeepSeek-V3
- **向量库**: VikingDB
- **可观测性**: Prometheus + Loki
- **部署**: Kubernetes

---

## 快速开始

```bash
# 启动本地环境
docker-compose up -d

# 执行数据库迁移
make migrate-up

# 启动API服务
make run-api
```

详见 [18_QuickStart_快速开始.md](./docs/18_QuickStart_快速开始.md)

---

## 文档目录

### 产品与设计
1. [01_Spec_产品需求.md](./docs/01_Spec_产品需求.md) - 用户故事、验收标准、功能优先级
2. [02_Plan_技术设计.md](./docs/02_Plan_技术设计.md) - 架构设计、Eino编排、多智能体协作
3. [03_API_接口契约.md](./docs/03_API_接口契约.md) - API定义、认证、错误码

### 部署与运维
4. [04_Deployment_部署手册.md](./docs/04_Deployment_部署手册.md) - 部署流程、配置、回滚
5. [05_Runbook_运行手册.md](./docs/05_Runbook_运行手册.md) - 值班SOP、故障处置、升级
6. [08_SRE_SLO与告警.md](./docs/08_SRE_SLO与告警.md) - SLO指标、告警分级
7. [17_Troubleshooting_故障排查.md](./docs/17_Troubleshooting_故障排查.md) - 常见问题排查

### 开发指南
8. [12_Dev_开发指南.md](./docs/12_Dev_开发指南.md) - 项目结构、开发环境、代码规范
9. [13_Database_数据库设计.md](./docs/13_Database_数据库设计.md) - 表结构、索引、迁移
10. [14_Tools_工具开发指南.md](./docs/14_Tools_工具开发指南.md) - 工具接口、实现示例
11. [15_Git_工作流.md](./docs/15_Git_工作流.md) - 分支策略、提交规范
12. [16_CICD_配置.md](./docs/16_CICD_配置.md) - CI/CD流程、Docker、K8s

### 质量与治理
13. [06_Test_测试策略.md](./docs/06_Test_测试策略.md) - 测试分层、关键用例
14. [07_Security_安全合规.md](./docs/07_Security_安全合规.md) - 权限模型、审批门禁
15. [09_Agent_Prompt策略.md](./docs/09_Agent_Prompt策略.md) - Agent角色、Prompt规范
16. [10_Eval_评测标准.md](./docs/10_Eval_评测标准.md) - 评测指标、通过标准
17. [11_KB_知识库治理.md](./docs/11_KB_知识库治理.md) - 知识来源、检索质量

### 项目管理
18. [19_Milestones_里程碑规划.md](./docs/19_Milestones_里程碑规划.md) - 6个里程碑、13周交付计划
19. [CLAUDE.md](./CLAUDE.md) - Claude Code工作指南

---

## 推荐阅读路径

### 产品/项目经理
01_Spec → 02_Plan → 19_Milestones → 10_Eval

### 架构师/研发工程师
02_Plan → 12_Dev → 13_Database → 14_Tools → 03_API → 16_CICD

### SRE/运维值班
05_Runbook → 04_Deployment → 08_SRE → 17_Troubleshooting → 07_Security

### 新人快速上手
18_QuickStart → 12_Dev → 02_Plan → CLAUDE.md

---

## 核心架构

### 四层架构
1. **接口层**: Web UI、CLI、Alert Webhook、OpenAPI
2. **编排层**: Eino Graph Runtime、Router、Multi-Agent Controller
3. **能力层**: ChatModel、Tool Runtime、Retriever、Policy Guard、Audit
4. **数据层**: Prometheus、Loki、Kubernetes API、VikingDB、Incident Store

### 多智能体系统
- **SupervisorAgent**: 任务拆解、协调、汇总
- **NetworkAgent**: 网络连通性诊断
- **DBAgent**: 数据库性能分析
- **AppAgent**: 应用异常诊断
- **InfraAgent**: 基础设施问题定位

---

## 关键SLO

- 对话响应P95: < 5s
- 告警分诊P95: < 30s
- RCA Top3命中率: ≥ 75%
- 动作成功率: ≥ 50%
- 高危误执行: 0次

---

## 开发命令

```bash
# 代码格式化
make fmt

# 代码检查
make lint

# 运行测试
make test

# 构建镜像
make docker-build

# 本地部署
make dev-up
```

---

## 当前版本

**v1.2** - 2026-03-24

---

## 贡献指南

1. Fork本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交变更 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建Pull Request

详见 [15_Git_工作流.md](./docs/15_Git_工作流.md)

---

## 许可证

[MIT License](LICENSE)
