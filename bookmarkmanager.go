package neo4j_tracing

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type BookmarkManagerTracer struct {
	neo4j.BookmarkManager
	tracer trace.Tracer
}

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
