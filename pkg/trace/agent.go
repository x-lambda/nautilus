package trace

import (
	"context"
	"fmt"
	"strings"

	"nautilus/pkg/conf"
	"nautilus/pkg/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	kindJaeger = "jaeger"
	kindZipkin = "zipkin"
)

// tp 全局TracerProvider
var tp *sdktrace.TracerProvider

func init() {
	ctx := context.TODO()
	config := &Config{
		Name:     "miniapp",
		Endpoint: conf.Get("OTEL_AGENT_ENDPOINT"),
		Sampler:  conf.GetFloat64("OTEL_AGENT_SAMPLE"),
		Batcher:  conf.Get("OTEL_AGENT_BATCH"),
	}

	if err := startAgent(ctx, config); err != nil {
		log.Get(ctx).Errorf("[otel] init agent err: %v", err)
	}
}

// startAgent 创建一个 Tracer Provider
func startAgent(ctx context.Context, c *Config) error {
	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(c.Sampler))),
		sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(c.Name))),
	}

	if len(c.Endpoint) > 0 {
		exporter, err := createExporter(c.Batcher, c.Endpoint)
		if err != nil {
			return err
		}

		opts = append(opts, sdktrace.WithBatcher(exporter))
	}

	tp = sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Get(ctx).Errorf("[otel] error: %v", err)
	}))

	return nil
}

// createExporter 选择opentelemetry的后端，目前只支持jaeger/zipkin
func createExporter(batcher string, endpoint string) (exporter sdktrace.SpanExporter, err error) {
	switch batcher {
	case kindJaeger:
		var opt jaeger.EndpointOption
		if strings.HasPrefix(endpoint, "http") {
			// HTTP
			opt = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint))
		} else {
			// UDP
			agentConfig := strings.SplitN(endpoint, ":", 2)
			if len(agentConfig) == 2 {
				opt = jaeger.WithAgentEndpoint(jaeger.WithAgentHost(agentConfig[0]), jaeger.WithAgentPort(agentConfig[1]))
			} else {
				opt = jaeger.WithAgentEndpoint(jaeger.WithAgentHost(agentConfig[0]))
			}
		}
		return jaeger.New(opt)
	case kindZipkin:
		return zipkin.New(endpoint)
	default:
		return nil, fmt.Errorf("unsupport exporter: %s", batcher)
	}
}

// Stop tracer provider shutdown
func Stop() {
	if tp == nil {
		return
	}

	tp.Shutdown(context.TODO())
}

// GetTraceID 提取trace id
func GetTraceID(ctx context.Context) (traceID string) {
	traceID = "no-trace-id"

	if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
		traceID = span.TraceID().String()
	}

	return
}

// version is the current release version of the gin instrumentation.
func version() string {
	return "0.0.1"
}

// SemVersion is the semantic version to be supplied to tracer/meter creation.
func SemVersion() string {
	return "semver:" + version()
}
