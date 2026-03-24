# 运维AI助手 - API接口契约

> 版本: v1.0  
> 更新日期: 2026-03-24  
> 适用范围: MVP（P0/P1）

---

## 1. 通用约定

- 基础路径：`/api/v1`
- 认证方式：`Authorization: Bearer <token>`
- 请求类型：`Content-Type: application/json`
- 时间格式：ISO8601（UTC）
- 幂等键：写接口建议携带 `Idempotency-Key`

统一响应格式：

```json
{
  "request_id": "req_xxx",
  "code": 0,
  "message": "ok",
  "data": {}
}
```

错误响应格式：

```json
{
  "request_id": "req_xxx",
  "code": 40001,
  "message": "invalid parameter",
  "error": {
    "field": "service",
    "reason": "required"
  }
}
```

---

## 2. 对话诊断

### 2.1 POST `/chat`

用途：对话式诊断与建议输出。

请求体：

```json
{
  "session_id": "sess_001",
  "query": "为什么payment-service延迟升高",
  "environment": "prod",
  "time_range": "30m"
}
```

响应体（成功）：

```json
{
  "request_id": "req_001",
  "code": 0,
  "message": "ok",
  "data": {
    "status": "critical",
    "summary": "疑似数据库连接池耗尽",
    "hypotheses": [
      {"description": "连接池耗尽", "confidence": 0.85},
      {"description": "下游超时", "confidence": 0.1},
      {"description": "网络抖动", "confidence": 0.05}
    ],
    "evidence": [
      {"type": "metric", "source": "prometheus", "ref": "query://..."},
      {"type": "log", "source": "loki", "ref": "logql://..."}
    ],
    "suggested_actions": [
      {"action": "检查连接池上限", "risk": "low"},
      {"action": "重启异常pod", "risk": "medium"}
    ]
  }
}
```

---

## 3. 告警接入

### 3.1 POST `/alerts/webhook`

用途：接收 Alertmanager 告警并触发分诊。

请求体：

```json
{
  "alert_id": "alert_123",
  "rule_name": "high_latency",
  "severity": "critical",
  "service": "payment-service",
  "environment": "prod",
  "starts_at": "2026-03-24T10:00:00Z",
  "labels": {"namespace": "payment"},
  "annotations": {"summary": "P99 latency too high"}
}
```

响应体：

```json
{
  "request_id": "req_002",
  "code": 0,
  "message": "accepted",
  "data": {
    "incident_id": "inc_001",
    "triage_status": "analyzing"
  }
}
```

---

## 4. 事件与报告查询

### 4.1 GET `/incidents/{incident_id}`

用途：获取事件详情、状态、最近诊断摘要。

### 4.2 GET `/reports/{incident_id}`

用途：获取结构化诊断报告与证据引用。

---

## 5. 动作审批与执行

### 5.1 POST `/actions/{action_id}/approve`

用途：高危动作审批（ops_admin）。

请求体：

```json
{
  "approved": true,
  "reason": "确认维护窗口内执行"
}
```

### 5.2 POST `/actions/{action_id}/execute`

用途：执行修复动作或回滚动作。

请求体：

```json
{
  "dry_run": false
}
```

---

## 6. 错误码

| code | 含义 |
|------|------|
| 0 | 成功 |
| 40001 | 参数错误 |
| 40101 | 未认证 |
| 40301 | 权限不足 |
| 40901 | 幂等冲突 |
| 42901 | 触发限流 |
| 50001 | 内部错误 |
| 50301 | 依赖服务不可用 |

---

## 7. 变更记录

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0 | 2026-03-24 | 初始版本 |

