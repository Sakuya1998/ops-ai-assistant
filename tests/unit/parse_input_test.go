package nodes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInputNode_Invoke(t *testing.T) {
	node := &ParseInputNode{}

	input := &InputMessage{
		Query:       "为什么payment-service延迟升高",
		Environment: "prod",
		TimeRange:   "30m",
	}

	result, err := node.Invoke(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, "payment-service", result.Service)
	assert.Equal(t, "prod", result.Environment)
}
