package tools

import (
	"context"
	"github.com/your-org/ops-ai-assistant/internal/datasource/loki"
)

type QueryLogsTool struct {
	client *loki.Client
}

func NewQueryLogsTool(client *loki.Client) *QueryLogsTool {
	return &QueryLogsTool{client: client}
}

func (t *QueryLogsTool) Name() string {
	return "query_logs"
}

func (t *QueryLogsTool) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
	query := params["query"].(string)
	limit := 100
	if l, ok := params["limit"].(int); ok {
		limit = l
	}

	result, err := t.client.Query(ctx, query, limit)
	if err != nil {
		return &ToolResult{Success: false, Error: err.Error()}, nil
	}

	return &ToolResult{Success: true, Data: result}, nil
}
