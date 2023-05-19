package neo4j_tracing

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
)

var _ neo4j.SessionWithContext = (*SessionWithContextTracer)(nil)

type SessionWithContextTracer struct {
	neo4j.SessionWithContext
	tracer trace.Tracer
}

func (s *SessionWithContextTracer) BeginTransaction(ctx context.Context, configurers ...func(config *neo4j.TransactionConfig)) (neo4j.ExplicitTransaction, error) {
	spanCtx, span := s.tracer.Start(ctx, "Session.BeginTransaction", trace.WithSpanKind(trace.SpanKindClient))

	tx, err := s.SessionWithContext.BeginTransaction(spanCtx, configurers...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()

		return nil, err
	}

	return NewExplicitTransactionTracer(spanCtx, tx, span, s.tracer), nil
}

func (s *SessionWithContextTracer) ExecuteRead(ctx context.Context, work neo4j.ManagedTransactionWork, configurers ...func(config *neo4j.TransactionConfig)) (_ any, err error) {
	spanCtx, span := s.tracer.Start(ctx, spanName("ExecuteRead"), trace.WithSpanKind(trace.SpanKindClient))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return s.SessionWithContext.ExecuteRead(spanCtx, func(tx neo4j.ManagedTransaction) (any, error) {
		txTracing := NewManagedTransactionTracer(spanCtx, tx, s.tracer)

		return work(txTracing)
	}, configurers...)
}

func (s *SessionWithContextTracer) ExecuteWrite(ctx context.Context, work neo4j.ManagedTransactionWork, configurers ...func(config *neo4j.TransactionConfig)) (_ any, err error) {
	spanCtx, span := s.tracer.Start(ctx, spanName("ExecuteWrite"), trace.WithSpanKind(trace.SpanKindClient))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return s.SessionWithContext.ExecuteWrite(spanCtx, func(tx neo4j.ManagedTransaction) (any, error) {
		txTracing := NewManagedTransactionTracer(spanCtx, tx, s.tracer)

		return work(txTracing)
	}, configurers...)
}

func (s *SessionWithContextTracer) Run(ctx context.Context, cypher string, params map[string]any, configurers ...func(config *neo4j.TransactionConfig)) (_ neo4j.ResultWithContext, err error) {
	spanCtx, span := s.tracer.Start(ctx, spanName("Run"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(semconv.DBStatement(cypher), semconv.DBSystemNeo4j))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return s.SessionWithContext.Run(spanCtx, cypher, params, configurers...)
}
