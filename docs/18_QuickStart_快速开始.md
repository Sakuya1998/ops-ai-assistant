# 运维AI助手 - 快速开始

> 版本: v1.0
> 更新日期: 2026-03-24

---

## 1. 5分钟快速体验

### 1.1 启动本地环境

```bash
# 克隆代码
git clone <repository-url>
cd ops-ai-assistant

# 启动依赖服务
docker-compose up -d

# 等待服务就绪
sleep 10

# 执行数据库迁移
make migrate-up

# 配置环境变量
export DEEPSEEK_API_KEY="your-api-key"

# 启动API服务
make run-api
```

### 1.2 测试对话诊断

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-001",
    "query": "为什么payment-service延迟升高",
    "environment": "prod"
  }'
```

### 1.3 测试告警接入

```bash
curl -X POST http://localhost:8080/api/v1/alerts/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "alert_id": "alert-001",
    "rule_name": "high_latency",
    "severity": "critical",
    "service": "payment-service",
    "environment": "prod",
    "starts_at": "2026-03-24T10:00:00Z"
  }'
```

---

## 2. 开发第一个功能

### 2.1 创建功能分支

```bash
git checkout -b feature/my-first-tool
```

### 2.2 实现一个简单工具

创建文件 `internal/capability/tools/echo_tool.go`:

```go
package tools

import "context"

type EchoTool struct{}

func (t *EchoTool) Name() string {
    return "echo"
}

func (t *EchoTool) Description() string {
    return "回显输入内容"
}

func (t *EchoTool) RiskLevel() RiskLevel {
    return RiskLow
}

func (t *EchoTool) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
    message := params["message"].(string)
    return &ToolResult{
        Success: true,
        Data:    map[string]string{"echo": message},
    }, nil
}
```

### 2.3 注册工具

在 `internal/capability/tools/registry.go` 中注册:

```go
func InitRegistry() *Registry {
    r := NewRegistry()
    r.Register(&EchoTool{})
    // ... 其他工具
    return r
}
```

### 2.4 测试

```bash
# 运行测试
go test ./internal/capability/tools/...

# 启动服务验证
make run-api
```

---

## 3. 下一步

- 阅读 [12_Dev_开发指南.md](./12_Dev_开发指南.md) 了解详细开发流程
- 阅读 [02_Plan_技术设计.md](./02_Plan_技术设计.md) 了解架构设计
- 查看 [14_Tools_工具开发指南.md](./14_Tools_工具开发指南.md) 学习工具开发

---

## 4. 变更记录

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0 | 2026-03-24 | 初始版本 |
