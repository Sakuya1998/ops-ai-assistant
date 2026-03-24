package nodes

import "context"

type InputMessage struct {
	Query       string
	Environment string
	TimeRange   string
}

type ParsedInput struct {
	Service     string
	Environment string
	TimeRange   string
	MetricType  string
}

type ParseInputNode struct{}

func (n *ParseInputNode) Invoke(ctx context.Context, input *InputMessage) (*ParsedInput, error) {
	return &ParsedInput{
		Service:     "payment-service",
		Environment: input.Environment,
		TimeRange:   input.TimeRange,
		MetricType:  "latency",
	}, nil
}
