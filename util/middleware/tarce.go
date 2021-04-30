package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"nautilus/util/ctxkit"
	"nautilus/util/trace"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func NewTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先分配span
		var span opentracing.Span
		c.Request, span = startSpan(c.Request)
		defer span.Finish()

		// trace id
		c.Request = newRequestID(c.Request, c.Writer)
	}
}

// newRequestID 每个请求创建一个trace id
func newRequestID(req *http.Request, resp http.ResponseWriter) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, ctxkit.StartTimeKey, time.Now())

	traceID := trace.GetTraceID(ctx)
	resp.Header().Set("x-trace-id", traceID)
	ctx = ctxkit.WithTraceID(ctx, traceID)

	return req.WithContext(ctx)
}

// startSpan 每个请求启动一个span
func startSpan(req *http.Request) (spanReq *http.Request, span opentracing.Span) {
	operation := "ServerHTTP"
	ctx := req.Context()

	tracer := opentracing.GlobalTracer()
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	if spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier); err == nil {
		span = opentracing.StartSpan(operation, ext.RPCServerOption(spanCtx))
		ctx = opentracing.ContextWithSpan(ctx, span)
	} else {
		fmt.Println("new request is here")
		span, ctx = opentracing.StartSpanFromContext(ctx, operation)
	}

	ext.SpanKindRPCServer.Set(span)
	span.SetTag(string(ext.HTTPUrl), req.URL.Path)

	spanReq = req.WithContext(ctx)
	return
}
