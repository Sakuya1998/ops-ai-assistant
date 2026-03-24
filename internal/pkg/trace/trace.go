package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var tracer oteltrace.Tracer

func Init(serviceName string) error {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)

	tracer = tp.Tracer(serviceName)
	return nil
}

func StartSpan(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	return tracer.Start(ctx, name)
}

func IDFromContext(ctx context.Context) string {
	span := oteltrace.SpanFromContext(ctx)
	return span.SpanContext().TraceID().String()
}
