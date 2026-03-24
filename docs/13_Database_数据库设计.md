# 运维AI助手 - 数据库设计

> 版本: v1.0
> 更新日期: 2026-03-24
> 数据库: PostgreSQL 14+

---

## 1. 表结构设计

### 1.1 incidents（事件表）

```sql
CREATE TABLE incidents (
    id VARCHAR(64) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    severity VARCHAR(20) NOT NULL, -- critical, high, medium, low
    status VARCHAR(20) NOT NULL, -- open, analyzing, resolved, closed
    source_alert_id VARCHAR(128),
    service VARCHAR(128) NOT NULL,
    environment VARCHAR(32) NOT NULL,
    started_at TIMESTAMP NOT NULL,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_service_env (service, environment),
    INDEX idx_status (status),
    INDEX idx_started_at (started_at)
);
```

### 1.2 hypotheses（根因假设表）

```sql
CREATE TABLE hypotheses (
    id VARCHAR(64) PRIMARY KEY,
    incident_id VARCHAR(64) NOT NULL,
    domain VARCHAR(32) NOT NULL, -- network, database, application, infrastructure
    description TEXT NOT NULL,
    confidence DECIMAL(3,2) NOT NULL, -- 0.00 ~ 1.00
    rank INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (incident_id) REFERENCES incidents(id) ON DELETE CASCADE,
    INDEX idx_incident_id (incident_id),
    INDEX idx_confidence (confidence DESC)
);
```

### 1.3 evidences（证据表）

```sql
CREATE TABLE evidences (
    id VARCHAR(64) PRIMARY KEY,
    incident_id VARCHAR(64) NOT NULL,
    hypothesis_id VARCHAR(64),
    type VARCHAR(32) NOT NULL, -- metric, log, event, trace
    source VARCHAR(64) NOT NULL, -- prometheus, loki, kubernetes
    summary TEXT NOT NULL,
    reference TEXT, -- 查询语句或链接
    captured_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (incident_id) REFERENCES incidents(id) ON DELETE CASCADE,
    FOREIGN KEY (hypothesis_id) REFERENCES hypotheses(id) ON DELETE SET NULL,
    INDEX idx_incident_id (incident_id),
    INDEX idx_type (type)
);
```

### 1.4 actions（动作记录表）

```sql
CREATE TABLE actions (
    id VARCHAR(64) PRIMARY KEY,
    incident_id VARCHAR(64) NOT NULL,
    action_type VARCHAR(64) NOT NULL, -- restart_pod, scale_deployment, exec_runbook
    risk_level VARCHAR(20) NOT NULL, -- low, medium, high
    status VARCHAR(20) NOT NULL, -- pending, approved, executing, succeeded, failed, rolled_back
    params JSONB NOT NULL,
    requested_by VARCHAR(128) NOT NULL,
    approved_by VARCHAR(128),
    approved_at TIMESTAMP,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    result JSONB,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (incident_id) REFERENCES incidents(id) ON DELETE CASCADE,
    INDEX idx_incident_id (incident_id),
    INDEX idx_status (status),
    INDEX idx_risk_level (risk_level)
);
```

### 1.5 audit_logs（审计日志表）

```sql
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    trace_id VARCHAR(64) NOT NULL,
    incident_id VARCHAR(64),
    action_id VARCHAR(64),
    user_id VARCHAR(128) NOT NULL,
    operation VARCHAR(64) NOT NULL,
    resource_type VARCHAR(64) NOT NULL,
    resource_id VARCHAR(128),
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_trace_id (trace_id),
    INDEX idx_incident_id (incident_id),
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
);
```

### 1.6 chat_sessions（对话会话表）

```sql
CREATE TABLE chat_sessions (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(128) NOT NULL,
    incident_id VARCHAR(64),
    environment VARCHAR(32),
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (incident_id) REFERENCES incidents(id) ON DELETE SET NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_last_active_at (last_active_at)
);
```

### 1.7 chat_messages（对话消息表）

```sql
CREATE TABLE chat_messages (
    id BIGSERIAL PRIMARY KEY,
    session_id VARCHAR(64) NOT NULL,
    role VARCHAR(20) NOT NULL, -- user, assistant, system
    content TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES chat_sessions(id) ON DELETE CASCADE,
    INDEX idx_session_id (session_id),
    INDEX idx_created_at (created_at)
);
```

---

## 2. 索引策略

### 2.1 复合索引

```sql
-- 事件查询优化
CREATE INDEX idx_incidents_service_status_time
ON incidents(service, status, started_at DESC);

-- 审计日志查询优化
CREATE INDEX idx_audit_logs_user_time
ON audit_logs(user_id, created_at DESC);
```

### 2.2 部分索引

```sql
-- 仅索引未解决的事件
CREATE INDEX idx_incidents_open
ON incidents(started_at DESC)
WHERE status IN ('open', 'analyzing');
```

---

## 3. 数据保留策略

```sql
-- 定期清理历史数据（通过定时任务执行）

-- 删除1年前的已关闭事件
DELETE FROM incidents
WHERE status = 'closed'
AND resolved_at < NOW() - INTERVAL '1 year';

-- 归档3年前的审计日志到对象存储
-- 然后删除本地记录
DELETE FROM audit_logs
WHERE created_at < NOW() - INTERVAL '3 years';
```

---

## 4. 迁移管理

使用 golang-migrate 管理数据库版本：

```bash
# 创建迁移文件
migrate create -ext sql -dir migrations -seq create_incidents_table

# 执行迁移
migrate -path migrations -database "postgres://..." up

# 回滚迁移
migrate -path migrations -database "postgres://..." down 1
```

---

## 5. 变更记录

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0 | 2026-03-24 | 初始版本 |
