package middleware

import (
	"fmt"

	"nautilus/pkg/ctxkit"
	xtrace "nautilus/pkg/trace"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerKey = "otel-go-contrib-tracer"
)

// NewTraceID otel trace中间件
// doc: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/7e31ebe040306aee2c826972269f938f9f0e7c34/instrumentation/github.com/gin-gonic/gin/otelgin/gintrace.go#L41
func NewTraceID() gin.HandlerFunc {
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("open-telemetry-go-contrib", trace.WithInstrumentationVersion(xtrace.SemVersion()))

	return func(c *gin.Context) {
		c.Set(tracerKey, tracer)
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()

		ctx := otel.GetTextMapPropagator().Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))
		opts := []trace.SpanStartOption{
			trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", c.Request)...),
			trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(c.Request)...),
			trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest("ServerHTTP", c.FullPath(), c.Request)...),
			trace.WithSpanKind(trace.SpanKindServer),
		}

		spanName := c.FullPath()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		}

		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		// 提取trace-id注入带response header中
		traceID := xtrace.GetTraceID(ctx)
		ctx = ctxkit.WithTraceID(ctx, traceID)
		c.Writer.Header().Set("x-trace-id", traceID)

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		status := c.Writer.Status()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)

		if len(c.Errors) > 0 {
			errStr := c.Errors.String()
			span.RecordError(fmt.Errorf(errStr))
			span.SetStatus(codes.Error, errStr)
		}
	}
}
