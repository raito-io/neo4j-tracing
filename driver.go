package neo4j_tracing

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/auth"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Neo4jTracer wraps a neo4j.Tracer object so the calls can be traced with open telemetry distributed tracing
type Neo4jTracer struct {
	tracer trace.Tracer
}

// NewNeo4jTracer creates an object that will wrap neo4j drivers with a tracing object
func NewNeo4jTracer(opts ...Option) *Neo4jTracer {
	cfg := config{}
	for _, o := range opts {
		o.apply(&cfg)
	}

	if cfg.TraceProvider == nil {
		cfg.TraceProvider = otel.GetTracerProvider()
	}

	return &Neo4jTracer{
		tracer: cfg.TraceProvider.Tracer(tracerName),
	}
}

// NewDriverWithContext is the entry point to the neo4j driver to create an instance of a neo4j.DriverWithContext that is wrapped by a tracing object
// More information about the arguments can be found in the underlying neo4j driver call neo4j.NewDriverWithContext
func (t *Neo4jTracer) NewDriverWithContext(target string, auth auth.TokenManager, configurers ...func(config2 *neo4j.Config)) (_ neo4j.DriverWithContext, err error) { //nolint:staticcheck
	driver, err := neo4j.NewDriverWithContext(target, auth, configurers...)
	if err != nil {
		return nil, err
	}

	return &DriverWithContextTracer{
		DriverWithContext: driver,
		tracer:            t.tracer,
	}, nil
}

type DriverWithContextTracer struct {
	neo4j.DriverWithContext
	tracer trace.Tracer
}

// ExecuteQueryBookmarkManager calls neo4j.DriverWithContext.ExecuteQueryBookmarkManager and wraps the resulting neo4j.BookmarkManager with a tracing object
func (n *DriverWithContextTracer) ExecuteQueryBookmarkManager() neo4j.BookmarkManager {
	return &BookmarkManagerTracer{
		BookmarkManager: n.DriverWithContext.ExecuteQueryBookmarkManager(),
		tracer:          n.tracer,
	}
}

// NewSession calls neo4j.DriverWithContext.NewSession and wraps the resulting neo4j.SessionWithContext with a tracing object
func (n *DriverWithContextTracer) NewSession(ctx context.Context, config neo4j.SessionConfig) neo4j.SessionWithContext {
	return &SessionWithContextTracer{
		SessionWithContext: n.DriverWithContext.NewSession(ctx, config),
		tracer:             n.tracer,
	}
}

// VerifyConnectivity calls neo4j.DriverWithContext.VerifyConnectivity and trace the call
func (n *DriverWithContextTracer) VerifyConnectivity(ctx context.Context) (err error) {
	spanCtx, span := n.tracer.Start(ctx, spanName("VerifyConnectivity"), trace.WithSpanKind(trace.SpanKindClient))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return n.DriverWithContext.VerifyConnectivity(spanCtx)
}

// VerifyAuthentication calls neo4j.DriverWithContext.VerifyAuthentication and trace the call
func (n *DriverWithContextTracer) VerifyAuthentication(ctx context.Context, auth *neo4j.AuthToken) (err error) {
	spanCtx, span := n.tracer.Start(ctx, spanName("VerifyAuthentication"), trace.WithSpanKind(trace.SpanKindClient))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return n.DriverWithContext.VerifyAuthentication(spanCtx, auth)
}

// GetServerInfo calls neo4j.GetServerInfo.VerifyConnectivity and trace the call
func (n *DriverWithContextTracer) GetServerInfo(ctx context.Context) (_ neo4j.ServerInfo, err error) {
	spanCtx, span := n.tracer.Start(ctx, spanName("GetServerInfo"), trace.WithSpanKind(trace.SpanKindClient))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return n.DriverWithContext.GetServerInfo(spanCtx)
}
