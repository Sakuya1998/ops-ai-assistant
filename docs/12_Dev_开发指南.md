# 运维AI助手 - 开发指南

> 版本: v1.0
> 更新日期: 2026-03-24
> 适用范围: 研发工程师

---

## 1. 技术栈

- **语言**: Go 1.21+
- **AI框架**: Eino (github.com/cloudwego/eino)
- **LLM**: DeepSeek-V3
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL 14+
- **缓存**: Redis 7+
- **向量库**: VikingDB
- **可观测性**: Prometheus + Loki + OpenTelemetry
- **容器化**: Docker + Kubernetes

---

## 2. 项目结构

```
ops-ai-assistant/
├── cmd/
│   ├── api/              # API服务入口
│   ├── worker/           # Worker服务入口
│   └── cli/              # CLI工具
├── internal/
│   ├── orchestration/    # Eino编排层
│   │   ├── graph/        # Graph定义
│   │   ├── nodes/        # 节点实现
│   │   └── router/       # 路由逻辑
│   ├── agents/           # 多智能体实现
│   │   ├── supervisor/
│   │   ├── network/
│   │   ├── database/
│   │   ├── application/
│   │   └── infrastructure/
│   ├── capability/       # 能力层
│   │   ├── llm/          # LLM封装
│   │   ├── tools/        # 工具实现
│   │   ├── retriever/    # RAG检索
│   │   ├── policy/       # 策略门禁
│   │   └── audit/        # 审计服务
│   ├── datasource/       # 数据源层
│   │   ├── prometheus/
│   │   ├── loki/
│   │   ├── kubernetes/
│   │   └── vikingdb/
│   ├── domain/           # 领域模型
│   │   ├── incident/
│   │   ├── hypothesis/
│   │   ├── evidence/
│   │   └── action/
│   ├── api/              # API处理器
│   │   ├── handler/
│   │   ├── middleware/
│   │   └── dto/
│   └── pkg/              # 公共包
│       ├── config/
│       ├── logger/
│       ├── trace/
│       └── errors/
├── pkg/                  # 可导出公共库
├── deployments/          # 部署配置
│   ├── kubernetes/
│   └── docker/
├── scripts/              # 脚本工具
├── tests/                # 测试
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── docs/                 # 文档
├── go.mod
└── Makefile
```

---

## 3. 开发环境搭建

### 3.1 前置依赖

```bash
# 安装Go 1.21+
brew install go

# 安装Docker Desktop
brew install --cask docker

# 安装kubectl
brew install kubectl

# 安装开发工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

### 3.2 本地开发环境

```bash
# 克隆代码
git clone <repository-url>
cd ops-ai-assistant

# 安装依赖
go mod download

# 启动本地依赖（Docker Compose）
make dev-up

# 运行数据库迁移
make migrate-up

# 启动API服务
make run-api

# 启动Worker服务
make run-worker
```

### 3.3 配置文件

创建 `config/local.yaml`:

```yaml
app:
  env: local
  log_level: debug

server:
  http_port: 8080

llm:
  provider: deepseek
  api_key: ${DEEPSEEK_API_KEY}
  model: deepseek-chat

datasources:
  prometheus:
    url: http://localhost:9090
  loki:
    url: http://localhost:3100
  kubernetes:
    in_cluster: false
    kubeconfig: ~/.kube/config

database:
  postgres:
    dsn: postgres://user:pass@localhost:5432/ops_ai?sslmode=disable
  redis:
    addr: localhost:6379

vikingdb:
  endpoint: http://localhost:8081
  collection: ops_knowledge
```

---

## 4. 核心模块开发指南

### 4.1 Eino Graph 编排

**节点定义示例**:

```go
// internal/orchestration/nodes/parse_input.go
package nodes

import (
    "context"
    "github.com/cloudwego/eino/compose"
)

type ParseInputNode struct {
    // 依赖注入
}

func (n *ParseInputNode) Invoke(ctx context.Context, input *InputMessage) (*ParsedInput, error) {
    // 解析告警或用户意图
    // 提取实体：服务名、环境、时间窗
    return &ParsedInput{
        Type: input.Type,
        Entities: extractEntities(input),
    }, nil
}
```

**Graph 构建示例**:

```go
// internal/orchestration/graph/diagnostic_graph.go
package graph

func NewDiagnosticGraph() *compose.Graph {
    g := compose.NewGraph[InputMessage, DiagnosticResult]()

    g.AddNode("parse_input", &nodes.ParseInputNode{})
    g.AddBranch("route_domain", &nodes.RouteDomainBranch{})
    g.AddNode("collect_metrics", &nodes.CollectMetricsNode{})
    g.AddNode("rca_rank", &nodes.RCARankNode{})
    g.AddNode("publish_report", &nodes.PublishReportNode{})

    g.AddEdge(compose.START, "parse_input")
    g.AddEdge("parse_input", "route_domain")
    g.AddEdge("route_domain", "collect_metrics")
    g.AddEdge("collect_metrics", "rca_rank")
    g.AddEdge("rca_rank", "publish_report")
    g.AddEdge("publish_report", compose.END)

    return g
}
```

### 4.2 Agent 实现

**Agent 接口定义**:

```go
// internal/agents/agent.go
package agents

type Agent interface {
    Name() string
    Diagnose(ctx context.Context, task *DiagnosticTask) (*AgentResult, error)
}

type AgentResult struct {
    Conclusion  string
    Confidence  float64
    Evidence    []Evidence
    Suggestions []string
}
```

**Specialist Agent 示例**:

```go
// internal/agents/network/network_agent.go
package network

type NetworkAgent struct {
    k8sClient    kubernetes.Interface
    promClient   prometheus.Client
    llm          llm.ChatModel
}

func (a *NetworkAgent) Diagnose(ctx context.Context, task *DiagnosticTask) (*AgentResult, error) {
    // 1. 收集网络相关证据
    evidence := a.collectNetworkEvidence(ctx, task)

    // 2. 调用LLM分析
    result := a.analyzeWithLLM(ctx, evidence)

    // 3. 返回结构化结果
    return &AgentResult{
        Conclusion: result.Conclusion,
        Confidence: result.Confidence,
        Evidence:   evidence,
    }, nil
}
```

### 4.3 工具实现

**工具接口**:

```go
// internal/capability/tools/tool.go
package tools

type Tool interface {
    Name() string
    Description() string
    RiskLevel() RiskLevel
    Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error)
}

type RiskLevel string

const (
    RiskLow    RiskLevel = "low"
    RiskMedium RiskLevel = "medium"
    RiskHigh   RiskLevel = "high"
)
```

**工具实现示例**:

```go
// internal/capability/tools/query_metrics.go
package tools

type QueryMetricsTool struct {
    promClient prometheus.Client
}

func (t *QueryMetricsTool) Name() string {
    return "query_metrics"
}

func (t *QueryMetricsTool) RiskLevel() RiskLevel {
    return RiskLow
}

func (t *QueryMetricsTool) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
    query := params["query"].(string)
    timeRange := params["time_range"].(string)

    result, err := t.promClient.Query(ctx, query, timeRange)
    if err != nil {
        return nil, err
    }

    return &ToolResult{
        Success: true,
        Data:    result,
    }, nil
}
```

### 4.4 策略门禁

```go
// internal/capability/policy/guard.go
package policy

type PolicyGuard struct {
    rbac RBACService
}

func (g *PolicyGuard) Check(ctx context.Context, action *Action) error {
    user := auth.UserFromContext(ctx)

    // RBAC检查
    if !g.rbac.HasPermission(user, action.Type) {
        return ErrPermissionDenied
    }

    // 风险级别检查
    if action.RiskLevel == RiskHigh && !action.Approved {
        return ErrApprovalRequired
    }

    return nil
}
```

---

## 5. 代码规范

### 5.1 命名规范

- 包名：小写单词，简短有意义
- 接口：名词或动词+er（如 `Agent`, `Retriever`）
- 结构体：大驼峰（如 `NetworkAgent`）
- 方法/函数：大驼峰（导出）或小驼峰（私有）
- 常量：大驼峰或全大写下划线分隔

### 5.2 错误处理

```go
// 使用自定义错误类型
var (
    ErrNotFound         = errors.New("resource not found")
    ErrPermissionDenied = errors.New("permission denied")
)

// 错误包装
if err != nil {
    return fmt.Errorf("failed to query metrics: %w", err)
}
```

### 5.3 日志规范

```go
// 使用结构化日志
logger.Info(ctx, "diagnostic started",
    "incident_id", incidentID,
    "service", service,
    "environment", env,
)

logger.Error(ctx, "tool execution failed",
    "tool", toolName,
    "error", err,
)
```

### 5.4 上下文传递

```go
// 必须传递context
func (s *Service) Process(ctx context.Context, input *Input) error {
    // 从context获取trace_id
    traceID := trace.IDFromContext(ctx)

    // 传递context到下游
    result, err := s.downstream.Call(ctx, input)

    return err
}
```

---

## 6. 测试规范

### 6.1 单元测试

```go
// internal/orchestration/nodes/parse_input_test.go
package nodes

func TestParseInputNode_Invoke(t *testing.T) {
    node := &ParseInputNode{}

    input := &InputMessage{
        Type: "alert",
        Content: "high latency in payment-service",
    }

    result, err := node.Invoke(context.Background(), input)

    assert.NoError(t, err)
    assert.Equal(t, "payment-service", result.Entities.Service)
}
```

### 6.2 集成测试

```go
// tests/integration/api_test.go
package integration

func TestChatAPI(t *testing.T) {
    // 启动测试服务器
    server := setupTestServer(t)
    defer server.Close()

    // 发送请求
    resp := sendChatRequest(server.URL, &ChatRequest{
        Query: "为什么payment-service延迟升高",
    })

    // 验证响应
    assert.Equal(t, 200, resp.StatusCode)
    assert.NotEmpty(t, resp.Data.Summary)
}
```

---

## 7. 常用命令

```bash
# 代码格式化
make fmt

# 代码检查
make lint

# 运行单元测试
make test

# 运行集成测试
make test-integration

# 生成API文档
make swagger

# 构建镜像
make docker-build

# 本地部署
make deploy-local
```

---

## 8. 变更记录

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0 | 2026-03-24 | 初始版本 |
