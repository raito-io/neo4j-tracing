package neo4j_tracing

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/auth"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Neo4jTracer struct {
	tracer trace.Tracer
}

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

func (n *DriverWithContextTracer) ExecuteQueryBookmarkManager() neo4j.BookmarkManager {
	return &BookmarkManagerTracer{
		BookmarkManager: n.DriverWithContext.ExecuteQueryBookmarkManager(),
		tracer:          n.tracer,
	}
}

func (n *DriverWithContextTracer) NewSession(ctx context.Context, config neo4j.SessionConfig) neo4j.SessionWithContext {
	return &SessionWithContextTracer{
		SessionWithContext: n.DriverWithContext.NewSession(ctx, config),
		tracer:             n.tracer,
	}
}

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
