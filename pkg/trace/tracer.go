package trace

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// TraceIDKey trace header
// https://www.w3.org/TR/trace-context/#trace-id
var TraceIDKey = http.CanonicalHeaderKey("x-trace-id")

type metadataSupplier struct {
	metadata *metadata.MD
}

func (m *metadataSupplier) Get(key string) string {
	values := m.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (m *metadataSupplier) Set(key string, value string) {
	m.metadata.Set(key, value)
}

func (m *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(*m.metadata))

	for key := range *m.metadata {
		out = append(out, key)
	}

	return out
}

var _ propagation.TextMapCarrier = new(metadataSupplier)

func Inject(ctx context.Context, p propagation.TextMapPropagator, metadata *metadata.MD) {
	p.Inject(ctx, &metadataSupplier{metadata: metadata})
}

func Extract(ctx context.Context, p propagation.TextMapPropagator, metadata *metadata.MD) (baggage.Baggage, trace.SpanContext) {
	ctx = p.Extract(ctx, &metadataSupplier{metadata: metadata})

	return baggage.FromContext(ctx), trace.SpanContextFromContext(ctx)
}

// InjectHeader 注入 open tracing 头信息
func InjectHeader(ctx trace.SpanContext, req *http.Request) {
	//opentracing.GlobalTracer().Inject(
	//	ctx,
	//	opentracing.HTTPHeaders,
	//	opentracing.HTTPHeadersCarrier(req.Header),
	//)
	//
	//jctx, ok := ctx.(jaeger.SpanContext)
	//if !ok {
	//	return
	//}

	//// Envoy 使用 Zipkin 风格头信息
	//// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/observability/tracing
	//req.Header.Set("x-b3-traceid", jctx.TraceID().String())
	//req.Header.Set("x-b3-spanid", jctx.SpanID().String())
	//req.Header.Set("x-b3-parentspanid", jctx.ParentID().String())
	//if jctx.IsSampled() {
	//	req.Header.Set("x-b3-sampled", "1")
	//}
	//if jctx.IsDebug() {
	//	req.Header.Set("x-b3-flags", "1")
	//}
}
