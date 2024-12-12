package neo4j_tracing

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
)

// ManagedTransactionTracer wraps a neo4j.ManagedTransaction object so the calls can be traced with open telemetry distributed tracing
type ManagedTransactionTracer struct {
	neo4j.ManagedTransaction
	ctx    context.Context
	tracer trace.Tracer
}

// NewManagedTransactionTracer returns a new ManagedTransactionTracer that wraps a neo4j.ManagedTransaction with correct tracing details
func NewManagedTransactionTracer(ctx context.Context, tx neo4j.ManagedTransaction, tracer trace.Tracer) *ManagedTransactionTracer {
	return &ManagedTransactionTracer{
		ManagedTransaction: tx,
		ctx:                ctx,
		tracer:             tracer,
	}
}

// Run calls neo4j.ManagedTransaction.Run and trace the call
func (t *ManagedTransactionTracer) Run(ctx context.Context, cypher string, params map[string]any) (_ neo4j.ResultWithContext, err error) {
	spanCtx, span := t.tracer.Start(t.ctx, spanName("Run"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(semconv.DBStatement(cypher), semconv.DBSystemNeo4j))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	result, err := t.ManagedTransaction.Run(ctx, cypher, params)

	return NewResultWithContextTracer(spanCtx, result, t.tracer), err
}

// ExplicitTransactionTracer wraps a neo4j.ExplicitTransaction object so the calls can be traced with open telemetry distributed tracing
type ExplicitTransactionTracer struct {
	neo4j.ExplicitTransaction
	ctx    context.Context
	txSpan trace.Span
	tracer trace.Tracer
}

// NewExplicitTransactionTracer returns a new ExplicitTransactionTracer that wraps a neo4j.ExplicitTransaction with correct tracing details
func NewExplicitTransactionTracer(ctx context.Context, tx neo4j.ExplicitTransaction, txSpan trace.Span, tracer trace.Tracer) *ExplicitTransactionTracer {
	return &ExplicitTransactionTracer{
		ExplicitTransaction: tx,
		ctx:                 ctx,
		txSpan:              txSpan,
		tracer:              tracer,
	}
}

// Run calls neo4j.ExplicitTransaction.Run and trace the call
func (t *ExplicitTransactionTracer) Run(ctx context.Context, cypher string, params map[string]any) (_ neo4j.ResultWithContext, err error) {
	spanCtx, span := t.tracer.Start(t.ctx, spanName("Run"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(semconv.DBStatement(cypher), semconv.DBSystemNeo4j))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	result, err := t.ExplicitTransaction.Run(ctx, cypher, params)

	return NewResultWithContextTracer(spanCtx, result, t.tracer), err
}

// Commit calls neo4j.ExplicitTransaction.Commit and trace the call
func (t *ExplicitTransactionTracer) Commit(ctx context.Context) (err error) {
	_, span := t.tracer.Start(t.ctx, spanName("Commit"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(semconv.DBSystemNeo4j))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return t.ExplicitTransaction.Commit(ctx)
}

// Rollback calls neo4j.ExplicitTransaction.Rollback and trace the call
func (t *ExplicitTransactionTracer) Rollback(ctx context.Context) (err error) {
	_, span := t.tracer.Start(t.ctx, spanName("Rollback"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(semconv.DBSystemNeo4j))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return t.ExplicitTransaction.Rollback(ctx)
}

// Close calls neo4j.ExplicitTransaction.Close and trace the call
func (t *ExplicitTransactionTracer) Close(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			t.txSpan.RecordError(err)
			t.txSpan.SetStatus(codes.Error, err.Error())
		}

		t.txSpan.End()
	}()

	return t.ExplicitTransaction.Close(ctx)
}
