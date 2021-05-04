package trace

import (
	"context"
	"io"
	"net/http"

	"nautilus/util/conf"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

// refer:
//	https://medium.com/opentracing/take-opentracing-for-a-hotrod-ride-f6e3141f7941
//	https://medium.com/opentracing/tracing-http-request-latency-in-go-with-opentracing-7cc1282a100a
//
// 开发环境可以一键部署jaeger:
// docker run -d -p6831:6831/udp -p16686:16686 jaegertracing/all-in-one:latest
// 然后配置
// 		JAEGER_TRACE_STATUS = 1
// 		JAEGER_TRACE_AGENT = "127.0.0.1:6831"
// 即可

var closer io.Closer

func init() {
	if conf.GetInt32("JAEGER_TRACE_STATUS") != 1 {
		return
	}

	var reporter jaeger.Reporter
	agent := conf.Get("JAEGER_TRACE_AGENT") // host+":"+port
	if agent == "" {
		reporter = jaeger.NewNullReporter()
	} else {
		// Jaeger tracer can be initialized with a transport that will
		// report tracing Spans to a Zipkin backend
		transport, _ := jaeger.NewUDPTransport(agent, 0)
		reporter = jaeger.NewRemoteReporter(transport)
	}

	serviceName := conf.Get("APP_ID")
	param := 0.9

	sampler, _ := jaeger.NewProbabilisticSampler(param)
	propagetor := zipkin.NewZipkinB3HTTPHeaderPropagator()
	tracer, c := jaeger.NewTracer(
		serviceName,
		sampler,
		reporter,
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, propagetor),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, propagetor),
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
	)

	closer = c
	opentracing.SetGlobalTracer(tracer)
}

// GetTraceID 从opentracing span中获取trace id
// https://github.com/opentracing/opentracing-go/issues/188
func GetTraceID(ctx context.Context) (traceID string) {
	traceID = "no-trace-id"

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}

	sc, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return
	}

	traceID = sc.TraceID().String()
	return
}

// InjectHeader 注入 open tracing 头信息
func InjectHeader(ctx opentracing.SpanContext, req *http.Request) {
	opentracing.GlobalTracer().Inject(
		ctx,
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	jctx, ok := ctx.(jaeger.SpanContext)
	if !ok {
		return
	}

	// TODO
	req.Header["Bili-Trace-Id"] = req.Header["Uber-Trace-Id"]

	// Envoy 使用 Zipkin 风格头信息
	// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/observability/tracing
	req.Header.Set("x-b3-traceid", jctx.TraceID().String())
	req.Header.Set("x-b3-spanid", jctx.SpanID().String())
	req.Header.Set("x-b3-parentspanid", jctx.ParentID().String())
	if jctx.IsSampled() {
		req.Header.Set("x-b3-sampled", "1")
	}
	if jctx.IsDebug() {
		req.Header.Set("x-b3-flags", "1")
	}
}

// StartFollowSpanFromContext 开启一个follow类型span
// follow类型用于异步任务，可能在root span结束之后才完成
func StartFollowSpanFromContext(ctx context.Context, operation string) (opentracing.Span, context.Context) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return opentracing.StartSpanFromContext(ctx, operation)
	}

	return opentracing.StartSpanFromContext(ctx, operation, opentracing.FollowsFrom(span.Context()))
}

// Stop 停止 trace 协程
func Stop() {
	closer.Close()
}
