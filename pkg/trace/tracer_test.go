package trace

import (
	"context"
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

const (
	traceIDStr = "4bf92f3577b34da6a3ce929d0e0e4736"
	spanIDStr  = "00f067aa0ba902b7"
)

var (
	traceID = mustTraceIDFromHex(traceIDStr)
	spanID  = mustSpanIDFromHex(spanIDStr)
)

func mustTraceIDFromHex(s string) (traceID trace.TraceID) {
	var err error
	traceID, err = trace.TraceIDFromHex(s)
	if err != nil {
		panic(err)
	}

	return
}

func mustSpanIDFromHex(s string) (spanID trace.SpanID) {
	var err error
	spanID, err = trace.SpanIDFromHex(s)
	if err != nil {
		panic(err)
	}

	return
}

func TestExtractValidTraceContext(t *testing.T) {
	stateStr := "key1=value1,key2=value2"
	state, err := trace.ParseTraceState(stateStr)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		traceparent string
		tracestate  string
		sc          trace.SpanContext
	}{
		{
			name:        "not sampled",
			traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "sampled",
			traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			}),
		},
		{
			name:        "valid tracestate",
			traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			tracestate:  stateStr,
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceState: state,
				Remote:     true,
			}),
		},
		{
			name:        "invalid tracestate perserves traceparent",
			traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			tracestate:  "invalid$@=invalid",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "future version not sampled",
			traceparent: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "future version sampled",
			traceparent: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			}),
		},
		{
			name:        "future version sample bit set",
			traceparent: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-09",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			}),
		},
		{
			name:        "future version sample bit not set",
			traceparent: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-08",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "future version addition data",
			traceparent: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00-XYZxsf09",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "B3 format ending in dash",
			traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00-",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "future version B3 format ending in dash",
			traceparent: "03-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00-",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	propagator := otel.GetTextMapPropagator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			md := metadata.MD{}
			md.Set("traceparent", tt.traceparent)
			md.Set("tracestate", tt.tracestate)
			_, spanCtx := Extract(ctx, propagator, &md)
			fmt.Println(tt.sc.TraceID().String())
			assert.Equal(t, tt.sc, spanCtx)
		})
	}
}
