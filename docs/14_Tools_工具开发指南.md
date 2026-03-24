# 运维AI助手 - 工具开发指南

> 版本: v1.0
> 更新日期: 2026-03-24

---

## 1. 工具接口定义

所有工具必须实现以下接口：

```go
package tools

import "context"

type Tool interface {
    Name() string
    Description() string
    RiskLevel() RiskLevel
    Schema() *ToolSchema
    Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error)
}

type RiskLevel string

const (
    RiskLow    RiskLevel = "low"
    RiskMedium RiskLevel = "medium"
    RiskHigh   RiskLevel = "high"
)

type ToolSchema struct {
    Parameters []Parameter
}

type Parameter struct {
    Name        string
    Type        string
    Required    bool
    Description string
}

type ToolResult struct {
    Success bool
    Data    interface{}
    Error   string
}
```

---

## 2. 工具实现示例

### 2.1 查询指标工具（低风险）

```go
package tools

import (
    "context"
    "github.com/your-org/ops-ai-assistant/internal/datasource/prometheus"
)

type QueryMetricsTool struct {
    promClient prometheus.Client
}

func NewQueryMetricsTool(client prometheus.Client) *QueryMetricsTool {
    return &QueryMetricsTool{promClient: client}
}

func (t *QueryMetricsTool) Name() string {
    return "query_metrics"
}

func (t *QueryMetricsTool) Description() string {
    return "查询Prometheus指标数据"
}

func (t *QueryMetricsTool) RiskLevel() RiskLevel {
    return RiskLow
}

func (t *QueryMetricsTool) Schema() *ToolSchema {
    return &ToolSchema{
        Parameters: []Parameter{
            {
                Name:        "query",
                Type:        "string",
                Required:    true,
                Description: "PromQL查询语句",
            },
            {
                Name:        "time_range",
                Type:        "string",
                Required:    false,
                Description: "时间范围，如30m, 1h",
            },
        },
    }
}

func (t *QueryMetricsTool) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
    query, ok := params["query"].(string)
    if !ok {
        return &ToolResult{
            Success: false,
            Error:   "query参数必须为字符串",
        }, nil
    }

    timeRange := "30m"
    if tr, ok := params["time_range"].(string); ok {
        timeRange = tr
    }

    result, err := t.promClient.Query(ctx, query, timeRange)
    if err != nil {
        return &ToolResult{
            Success: false,
            Error:   err.Error(),
        }, nil
    }

    return &ToolResult{
        Success: true,
        Data:    result,
    }, nil
}
```

### 2.2 重启Pod工具（中风险）

```go
package tools

import (
    "context"
    "k8s.io/client-go/kubernetes"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RestartPodTool struct {
    k8sClient kubernetes.Interface
}

func (t *RestartPodTool) Name() string {
    return "restart_pod"
}

func (t *RestartPodTool) RiskLevel() RiskLevel {
    return RiskMedium
}

func (t *RestartPodTool) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
    namespace := params["namespace"].(string)
    podName := params["pod_name"].(string)

    // 删除Pod，由控制器重建
    err := t.k8sClient.CoreV1().Pods(namespace).Delete(
        ctx,
        podName,
        metav1.DeleteOptions{},
    )

    if err != nil {
        return &ToolResult{
            Success: false,
            Error:   err.Error(),
        }, nil
    }

    return &ToolResult{
        Success: true,
        Data: map[string]string{
            "namespace": namespace,
            "pod":       podName,
            "action":    "deleted",
        },
    }, nil
}
```

---

## 3. 工具注册

```go
package tools

type Registry struct {
    tools map[string]Tool
}

func NewRegistry() *Registry {
    return &Registry{
        tools: make(map[string]Tool),
    }
}

func (r *Registry) Register(tool Tool) {
    r.tools[tool.Name()] = tool
}

func (r *Registry) Get(name string) (Tool, bool) {
    tool, ok := r.tools[name]
    return tool, ok
}

func (r *Registry) List() []Tool {
    tools := make([]Tool, 0, len(r.tools))
    for _, tool := range r.tools {
        tools = append(tools, tool)
    }
    return tools
}
```

---

## 4. 工具清单

| 工具名 | 风险级别 | 说明 |
|-------|---------|------|
| query_metrics | 低 | 查询Prometheus指标 |
| query_logs | 低 | 查询Loki日志 |
| get_k8s_resource | 低 | 获取K8s资源描述 |
| get_k8s_events | 低 | 获取K8s事件 |
| restart_pod | 中 | 重启单个Pod |
| scale_deployment | 高 | 调整副本数 |
| exec_runbook | 中/高 | 执行修复剧本 |

---

## 5. 变更记录

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0 | 2026-03-24 | 初始版本 |
