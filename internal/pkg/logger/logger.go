package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
)

// Logger 接口抽象
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
	WithGroup(name string) Logger
}

type slogLogger struct {
	logger *slog.Logger
	level  *atomic.Value
}

// dynamicLevelHandler 支持动态级别的Handler
type dynamicLevelHandler struct {
	handler slog.Handler
	level   *atomic.Value
}

func (h *dynamicLevelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := h.level.Load().(slog.Level)
	return level >= minLevel
}

func (h *dynamicLevelHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

func (h *dynamicLevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &dynamicLevelHandler{handler: h.handler.WithAttrs(attrs), level: h.level}
}

func (h *dynamicLevelHandler) WithGroup(name string) slog.Handler {
	return &dynamicLevelHandler{handler: h.handler.WithGroup(name), level: h.level}
}

var defaultLogger *slogLogger

type contextKey string

const (
	traceIDKey   contextKey = "trace_id"
	requestIDKey contextKey = "request_id"
)

// Init 初始化logger
func Init(level string) error {
	logLevel := parseLevel(level)
	levelVar := &atomic.Value{}
	levelVar.Store(logLevel)

	opts := &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replaceAttr,
	}

	baseHandler := slog.NewJSONHandler(os.Stdout, opts)
	dynamicHandler := &dynamicLevelHandler{
		handler: baseHandler,
		level:   levelVar,
	}

	defaultLogger = &slogLogger{
		logger: slog.New(dynamicHandler),
		level:  levelVar,
	}

	slog.SetDefault(defaultLogger.logger)
	return nil
}

// SetLevel 动态更新日志级别
func SetLevel(level string) {
	if defaultLogger == nil {
		return
	}
	logLevel := parseLevel(level)
	defaultLogger.level.Store(logLevel)
}

// WithTraceID 添加trace_id到context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// WithRequestID 添加request_id到context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		if source, ok := a.Value.Any().(*slog.Source); ok {
			if strings.Contains(source.File, "internal/pkg/logger") {
				var pcs [1]uintptr
				runtime.Callers(5, pcs[:])
				fs := runtime.CallersFrames(pcs[:])
				f, _ := fs.Next()
				source.File = f.File
				source.Line = f.Line
			}
		}
	}
	return a
}

func extractContextFields(ctx context.Context) []any {
	traceID, hasTrace := ctx.Value(traceIDKey).(string)
	requestID, hasRequest := ctx.Value(requestIDKey).(string)

	if !hasTrace && !hasRequest {
		return nil
	}

	fields := make([]any, 0, 4)
	if hasTrace && traceID != "" {
		fields = append(fields, "trace_id", traceID)
	}
	if hasRequest && requestID != "" {
		fields = append(fields, "request_id", requestID)
	}
	return fields
}

// Logger接口实现
func (l *slogLogger) Debug(ctx context.Context, msg string, args ...any) {
	if fields := extractContextFields(ctx); fields != nil {
		args = append(fields, args...)
	}
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	if fields := extractContextFields(ctx); fields != nil {
		args = append(fields, args...)
	}
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *slogLogger) Warn(ctx context.Context, msg string, args ...any) {
	if fields := extractContextFields(ctx); fields != nil {
		args = append(fields, args...)
	}
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	if fields := extractContextFields(ctx); fields != nil {
		args = append(fields, args...)
	}
	l.logger.ErrorContext(ctx, msg, args...)
}

func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{logger: l.logger.With(args...), level: l.level}
}

func (l *slogLogger) WithGroup(name string) Logger {
	return &slogLogger{logger: l.logger.WithGroup(name), level: l.level}
}

// 包级便捷函数
func Debug(ctx context.Context, msg string, args ...any) {
	if defaultLogger == nil {
		slog.DebugContext(ctx, msg, args...)
		return
	}
	defaultLogger.Debug(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	if defaultLogger == nil {
		slog.InfoContext(ctx, msg, args...)
		return
	}
	defaultLogger.Info(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	if defaultLogger == nil {
		slog.WarnContext(ctx, msg, args...)
		return
	}
	defaultLogger.Warn(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	if defaultLogger == nil {
		slog.ErrorContext(ctx, msg, args...)
		return
	}
	defaultLogger.Error(ctx, msg, args...)
}

func With(args ...any) Logger {
	if defaultLogger == nil {
		return &slogLogger{logger: slog.Default().With(args...)}
	}
	return defaultLogger.With(args...)
}

func WithGroup(name string) Logger {
	if defaultLogger == nil {
		return &slogLogger{logger: slog.Default().WithGroup(name)}
	}
	return defaultLogger.WithGroup(name)
}
