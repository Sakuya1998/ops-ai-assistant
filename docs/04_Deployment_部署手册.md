# 运维AI助手 - 部署手册

> 版本: v1.0  
> 更新日期: 2026-03-24  
> 环境: Kubernetes

---

## 1. 部署目标

- 完成 API、Worker、Web 三个组件可用部署
- 接入 Prometheus、Loki、Kubernetes API、VikingDB
- 满足基础可观测与回滚能力

---

## 2. 前置条件

- Kubernetes 集群可用（建议 1.27+）
- 命名空间：`ops-ai`
- 已准备以下依赖：
  - Redis
  - PostgreSQL
  - VikingDB
  - Prometheus / Loki
- 可用镜像仓库与镜像拉取密钥

---

## 3. 配置项清单

| 变量名 | 说明 |
|-------|------|
| `APP_ENV` | 运行环境，如 prod/staging |
| `HTTP_PORT` | API 服务端口 |
| `DEEPSEEK_API_KEY` | 模型访问密钥 |
| `PROMETHEUS_URL` | 指标查询地址 |
| `LOKI_URL` | 日志查询地址 |
| `K8S_CLUSTER_NAME` | 集群标识 |
| `POSTGRES_DSN` | 元数据数据库连接串 |
| `REDIS_ADDR` | 缓存地址 |
| `VIKINGDB_ENDPOINT` | 向量库地址 |

---

## 4. 部署步骤

1. 创建命名空间与基础密钥  
2. 应用 ConfigMap 和 Secret  
3. 部署 `ops-ai-api`、`ops-ai-worker`、`ops-ai-web`  
4. 配置 Service 和 Ingress  
5. 执行健康检查与冒烟测试  

---

## 5. 健康检查

- API 存活检查：`/healthz`
- API 就绪检查：`/readyz`
- Worker 心跳指标：`ops_ai_worker_heartbeat`
- Web 可用性：首页 200

---

## 6. 冒烟测试

- 调用 `/api/v1/chat` 返回结构化响应
- 调用 `/api/v1/alerts/webhook` 返回 `incident_id`
- 查询 `/api/v1/incidents/{id}` 可返回状态

---

## 7. 回滚策略

- 使用 Deployment revision 回滚至上一稳定版本
- 恢复前一个 ConfigMap/Secret 版本
- 验证关键接口和SLO指标恢复正常

---

## 8. 上线检查清单

- [ ] 配置项齐全且密钥已注入
- [ ] RBAC 与网络策略已生效
- [ ] 监控、日志、追踪正常
- [ ] 高风险动作审批链路可用
- [ ] 冒烟测试通过

---

## 9. 变更记录

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0 | 2026-03-24 | 初始版本 |

