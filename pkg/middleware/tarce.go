package middleware

import (
	"net/http"

	"nautilus/pkg/ctxkit"
	"nautilus/pkg/trace"

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
	traceID := trace.GetTraceID(ctx)
	ctx = ctxkit.WithTraceID(ctx, traceID)
	resp.Header().Set("x-trace-id", traceID)

	return req.WithContext(ctx)
}

// startSpan 每个请求启动一个span
func startSpan(req *http.Request) (spanReq *http.Request, span opentracing.Span) {
	operation := "ServerHTTP"
	ctx := req.Context()

	tracer := opentracing.GlobalTracer()
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	// 这里会尝试从 http headers 提取出trace信息
	if spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier); err == nil {
		// 能提取出，说明上游有调用方，则在上游的span中续上span[child span]
		span = opentracing.StartSpan(operation, opentracing.ChildOf(spanCtx))
		ctx = opentracing.ContextWithSpan(ctx, span)
	} else {
		// 相当于创建一个 root span
		span, ctx = opentracing.StartSpanFromContext(ctx, operation)
	}

	span.SetTag(string(ext.HTTPUrl), req.URL.Path)

	spanReq = req.WithContext(ctx)
	return
}
