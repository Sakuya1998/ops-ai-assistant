package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryMetricsTool_Execute(t *testing.T) {
	tool := &QueryMetricsTool{}

	params := map[string]interface{}{
		"query":      "up",
		"time_range": "5m",
	}

	result, err := tool.Execute(context.Background(), params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}
