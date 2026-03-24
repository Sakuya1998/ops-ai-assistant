package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/your-org/ops-ai-assistant/internal/api/handler"
	"github.com/your-org/ops-ai-assistant/internal/pkg/config"
	"github.com/your-org/ops-ai-assistant/internal/pkg/logger"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.Init(cfg.App.LogLevel); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	chatHandler := handler.NewChatHandler()
	r.POST("/api/v1/chat", chatHandler.Handle)

	r.Run(":8080")
}
