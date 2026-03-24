package nodes

import (
	"context"
	"github.com/Sakuya1998/ops-ai-assistant/internal/capability/tools"
)

type Evidence struct {
	Type    string
	Source  string
	Summary string
}

type CollectMetricsNode struct {
	tool *tools.QueryMetricsTool
}

func NewCollectMetricsNode(tool *tools.QueryMetricsTool) *CollectMetricsNode {
	return &CollectMetricsNode{tool: tool}
}

func (n *CollectMetricsNode) Invoke(ctx context.Context, input *ParsedInput) ([]Evidence, error) {
	params := map[string]interface{}{
		"query":      "rate(http_request_duration_seconds_sum[5m])",
		"time_range": input.TimeRange,
	}

	result, _ := n.tool.Execute(ctx, params)

	return []Evidence{
		{Type: "metric", Source: "prometheus", Summary: "P99延迟升高"},
	}, nil
}
