package tools

import "context"

type Tool interface {
	Name() string
	Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error)
}

type ToolResult struct {
	Success bool
	Data    interface{}
	Error   string
}
