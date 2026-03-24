# 运维AI助手 - 故障排查手册

> 版本: v1.0
> 更新日期: 2026-03-24

---

## 1. 常见问题

### 1.1 API响应超时

**症状**: `/api/v1/chat` 接口响应时间超过30秒

**排查步骤**:
```bash
# 1. 检查API服务日志
kubectl logs -n ops-ai deployment/ops-ai-api --tail=100

# 2. 检查LLM调用延迟
# 查看Prometheus指标
ops_ai_llm_request_duration_seconds

# 3. 检查数据源连接
# Prometheus连接
curl http://prometheus:9090/-/healthy

# Loki连接
curl http://loki:3100/ready
```

**解决方案**:
- LLM超时: 增加超时配置或切换模型
- 数据源慢: 优化查询语句或增加缓存
- 并发过高: 增加副本数或限流

---

### 1.2 告警分诊失败

**症状**: Webhook接收告警但未生成Incident

**排查步骤**:
```bash
# 1. 检查Webhook日志
kubectl logs -n ops-ai deployment/ops-ai-api | grep webhook

# 2. 检查数据库连接
psql -h postgres -U ops_ai -d ops_ai -c "SELECT 1"

# 3. 检查Worker状态
kubectl get pods -n ops-ai -l app=ops-ai-worker
```

**解决方案**:
- 数据库连接失败: 检查连接串和网络
- Worker未运行: 重启Worker服务
- 告警格式错误: 检查Webhook请求体格式

---

### 1.3 高危动作未拦截

**症状**: 高风险动作未经审批直接执行

**排查步骤**:
```bash
# 1. 检查审计日志
SELECT * FROM audit_logs
WHERE operation = 'execute_action'
AND created_at > NOW() - INTERVAL '1 hour'
ORDER BY created_at DESC;

# 2. 检查策略配置
kubectl get configmap -n ops-ai ops-ai-config -o yaml | grep policy

# 3. 检查RBAC配置
kubectl get rolebinding -n ops-ai
```

**解决方案**:
- 策略未生效: 重启服务加载配置
- RBAC配置错误: 修正权限配置
- 代码Bug: 检查PolicyGuard实现

---

## 2. 性能问题

### 2.1 内存泄漏

**监控指标**:
```promql
# 内存使用趋势
container_memory_usage_bytes{namespace="ops-ai"}

# Goroutine数量
go_goroutines{job="ops-ai-api"}
```

**排查工具**:
```bash
# 获取pprof数据
curl http://ops-ai-api:8080/debug/pprof/heap > heap.prof

# 分析内存
go tool pprof heap.prof
```

---

### 2.2 数据库连接池耗尽

**症状**: 大量 "too many connections" 错误

**排查**:
```sql
-- 查看当前连接数
SELECT count(*) FROM pg_stat_activity;

-- 查看最大连接数
SHOW max_connections;
```

**解决**:
- 增加连接池大小
- 优化慢查询
- 添加连接超时

---

## 3. 数据问题

### 3.3 证据丢失

**症状**: Incident有假设但无证据

**排查**:
```sql
SELECT i.id, i.title,
       COUNT(h.id) as hypothesis_count,
       COUNT(e.id) as evidence_count
FROM incidents i
LEFT JOIN hypotheses h ON i.id = h.incident_id
LEFT JOIN evidences e ON i.id = e.incident_id
WHERE i.created_at > NOW() - INTERVAL '1 day'
GROUP BY i.id, i.title
HAVING COUNT(e.id) = 0;
```

**解决**:
- 检查数据源连接
- 检查证据收集节点日志
- 验证外键约束

---

## 4. 变更记录

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0 | 2026-03-24 | 初始版本 |
