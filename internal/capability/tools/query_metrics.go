package tools

import (
	"context"
	"github.com/Sakuya1998/ops-ai-assistant/internal/datasource/prometheus"
)

type QueryMetricsTool struct {
	client *prometheus.Client
}

func NewQueryMetricsTool(client *prometheus.Client) *QueryMetricsTool {
	return &QueryMetricsTool{client: client}
}

func (t *QueryMetricsTool) Name() string {
	return "query_metrics"
}

func (t *QueryMetricsTool) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
	query := params["query"].(string)
	timeRange := "30m"
	if tr, ok := params["time_range"].(string); ok {
		timeRange = tr
	}

	result, err := t.client.Query(ctx, query, timeRange)
	if err != nil {
		return &ToolResult{Success: false, Error: err.Error()}, nil
	}

	return &ToolResult{Success: true, Data: result}, nil
}
