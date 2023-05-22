package neo4j_tracing

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// BookmarkManagerTracer wraps a neo4j.BookmarkManager object so the calls can be traced with open telemetry distributed tracing
type BookmarkManagerTracer struct {
	// Actual neo4j BookmarkManager
	neo4j.BookmarkManager

	// OTEL tracer
	tracer trace.Tracer
}

// UpdateBookmarks calls neo4j.BookmarkManager.UpdateBookmarks and trace the call
func (b *BookmarkManagerTracer) UpdateBookmarks(ctx context.Context, previousBookmarks, newBookmarks neo4j.Bookmarks) (err error) {
	spanCtx, span := b.tracer.Start(ctx, spanName("UpdateBookmarks"), trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.StringSlice("previousBookmarks", previousBookmarks), attribute.StringSlice("newBookmarks", newBookmarks)))

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return b.BookmarkManager.UpdateBookmarks(spanCtx, previousBookmarks, newBookmarks)
}

// GetBookmarks calls neo4j.BookmarkManager.GetBookmarks and trace the call
func (b *BookmarkManagerTracer) GetBookmarks(ctx context.Context) (_ neo4j.Bookmarks, err error) {
	spanCtx, span := b.tracer.Start(ctx, spanName("GetBookmarks"), trace.WithSpanKind(trace.SpanKindClient))

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	return b.BookmarkManager.GetBookmarks(spanCtx)
}
