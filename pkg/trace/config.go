package trace

// TracerName tracing name
const TracerName = "miniapp"

// Config opentelemetry collect config
type Config struct {
	Name     string `json:"name,optional"`
	Endpoint string `json:"endpoint,optional"`

	// Sampler 采样比例
	Sampler float64 `json:"sampler,default=1.0"`

	// Batcher otel后端，只支持jaeger/zipkin
	// registers the exporter with the TracerProvider
	Batcher string `json:"batcher,default=jaeger,options=jaeger|zipkin"`
}
