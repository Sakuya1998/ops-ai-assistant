package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/ops-ai-assistant/internal/api/dto"
	"github.com/your-org/ops-ai-assistant/internal/orchestration/nodes"
)

type ChatHandler struct{}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

func (h *ChatHandler) Handle(c *gin.Context) {
	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.ChatResponse{
			RequestID: uuid.New().String(),
			Code:      40001,
			Message:   "invalid request",
		})
		return
	}

	input := &nodes.InputMessage{
		Query:       req.Query,
		Environment: req.Environment,
		TimeRange:   req.TimeRange,
	}

	parseNode := &nodes.ParseInputNode{}
	parsed, _ := parseNode.Invoke(c.Request.Context(), input)

	c.JSON(200, dto.ChatResponse{
		RequestID: uuid.New().String(),
		Code:      0,
		Message:   "ok",
		Data: &dto.DiagnosticResult{
			Status:  "warning",
			Summary: "payment-service延迟升高，疑似数据库连接池问题",
			Evidence: []dto.Evidence{
				{Type: "metric", Source: "prometheus", Summary: "P99延迟从50ms升至200ms"},
			},
		},
	})
}
