package neo4j_tracing

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
)

type ManagedTransactionTracer struct {
	neo4j.ManagedTransaction
	ctx    context.Context
	tracer trace.Tracer
}

func NewManagedTransactionTracer(ctx context.Context, tx neo4j.ManagedTransaction, tracer trace.Tracer) *ManagedTransactionTracer {
	return &ManagedTransactionTracer{
		ManagedTransaction: tx,
		ctx:                ctx,
		tracer:             tracer,
	}
}

func (t *ManagedTransactionTracer) Run(ctx context.Context, cypher string, params map[string]any) (_ neo4j.ResultWithContext, err error) {
	_, span := t.tracer.Start(t.ctx, spanName("Run"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(semconv.DBStatement(cypher), semconv.DBSystemNeo4j))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return t.ManagedTransaction.Run(ctx, cypher, params)
}

type ExplicitTransactionTracer struct {
	neo4j.ExplicitTransaction
	ctx    context.Context
	txSpan trace.Span
	tracer trace.Tracer
}

func NewExplicitTransactionTracer(ctx context.Context, tx neo4j.ExplicitTransaction, txSpan trace.Span, tracer trace.Tracer) *ExplicitTransactionTracer {
	return &ExplicitTransactionTracer{
		ExplicitTransaction: tx,
		ctx:                 ctx,
		txSpan:              txSpan,
		tracer:              tracer,
	}
}

func (t *ExplicitTransactionTracer) Run(ctx context.Context, cypher string, params map[string]any) (_ neo4j.ResultWithContext, err error) {
	_, span := t.tracer.Start(t.ctx, spanName("Run"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(semconv.DBStatement(cypher), semconv.DBSystemNeo4j))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return t.ExplicitTransaction.Run(ctx, cypher, params)
}

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
