package neo4j_tracing

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ResultWithContextTracer struct {
	neo4j.ResultWithContext
	ctx    context.Context
	tracer trace.Tracer
}

func NewResultWithContextTracer(ctx context.Context, result neo4j.ResultWithContext, tracer trace.Tracer) neo4j.ResultWithContext {
	return &ResultWithContextTracer{ResultWithContext: result, ctx: ctx, tracer: tracer}
}

func (r *ResultWithContextTracer) NextRecord(ctx context.Context, record **neo4j.Record) bool {
	_, span := r.tracer.Start(r.ctx, spanName("Record.NextRecord"), trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	return r.ResultWithContext.NextRecord(ctx, record)
}

func (r *ResultWithContextTracer) Next(ctx context.Context) bool {
	_, span := r.tracer.Start(r.ctx, spanName("Record.Next"), trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	return r.ResultWithContext.Next(ctx)
}

func (r *ResultWithContextTracer) PeekRecord(ctx context.Context, record **neo4j.Record) bool {
	_, span := r.tracer.Start(r.ctx, spanName("Record.PeekRecord"), trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	return r.ResultWithContext.PeekRecord(ctx, record)
}

func (r *ResultWithContextTracer) Peek(ctx context.Context) bool {
	_, span := r.tracer.Start(r.ctx, spanName("Record.Peek"), trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	return r.ResultWithContext.Peek(ctx)
}

func (r *ResultWithContextTracer) Collect(ctx context.Context) (_ []*neo4j.Record, err error) {
	_, span := r.tracer.Start(r.ctx, spanName("Record.Collect"), trace.WithSpanKind(trace.SpanKindInternal))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return r.ResultWithContext.Collect(ctx)
}

func (r *ResultWithContextTracer) Single(ctx context.Context) (_ *neo4j.Record, err error) {
	_, span := r.tracer.Start(r.ctx, spanName("Record.Single"), trace.WithSpanKind(trace.SpanKindInternal))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return r.ResultWithContext.Single(ctx)
}

func (r *ResultWithContextTracer) Consume(ctx context.Context) (_ neo4j.ResultSummary, err error) {
	_, span := r.tracer.Start(r.ctx, spanName("Record.Consume"), trace.WithSpanKind(trace.SpanKindInternal))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return r.ResultWithContext.Consume(ctx)
}
