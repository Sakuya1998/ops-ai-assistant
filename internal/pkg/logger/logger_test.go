package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	err := Init("debug")
	assert.NoError(t, err)
	assert.NotNil(t, defaultLogger)
}

func TestDynamicLevelChange(t *testing.T) {
	var buf bytes.Buffer
	levelVar := &atomic.Value{}
	levelVar.Store(slog.LevelInfo)

	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
	dynamicHandler := &dynamicLevelHandler{handler: baseHandler, level: levelVar}
	logger := slog.New(dynamicHandler)

	ctx := context.Background()

	// Info级别，debug不应输出
	logger.Debug("debug message")
	assert.NotContains(t, buf.String(), "debug message")

	// 切换到Debug级别
	levelVar.Store(slog.LevelDebug)
	logger.Debug("debug message 2")
	assert.Contains(t, buf.String(), "debug message 2")
}

func TestContextFields(t *testing.T) {
	var buf bytes.Buffer
	levelVar := &atomic.Value{}
	levelVar.Store(slog.LevelInfo)

	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
	dynamicHandler := &dynamicLevelHandler{handler: baseHandler, level: levelVar}
	defaultLogger = &slogLogger{logger: slog.New(dynamicHandler), level: levelVar}

	ctx := WithTraceID(context.Background(), "trace-123")
	ctx = WithRequestID(ctx, "req-456")

	Info(ctx, "test message")

	var logEntry map[string]interface{}
	json.Unmarshal(buf.Bytes(), &logEntry)

	assert.Equal(t, "trace-123", logEntry["trace_id"])
	assert.Equal(t, "req-456", logEntry["request_id"])
}

func TestSetLevel(t *testing.T) {
	Init("info")
	SetLevel("debug")

	level := defaultLogger.level.Load().(slog.Level)
	assert.Equal(t, slog.LevelDebug, level)
}

