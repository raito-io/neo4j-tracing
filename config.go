package neo4j_tracing

import "go.opentelemetry.io/otel/trace"

type config struct {
	TraceProvider trace.TracerProvider
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return optionFunc(func(c *config) {
		c.TraceProvider = tp
	})
}
