package pixel

import (
	"context"

	tracing "github.com/ricardo-ch/go-tracing"
)

type pxTracing struct {
	next Service
}

// NewpxTracing ...
func NewpxTracing(i Service) Service {
	return pxTracing{
		next: i,
	}
}

func (s pxTracing) Getpx(ctx context.Context, id string) (m *Px, err error) {
	span, ctx := tracing.CreateSpan(ctx, "pxs.service::Getpx", &map[string]interface{}{"pxID": id})
	defer func() {
		if err != nil {
			tracing.SetSpanError(span, err)
		}
		span.Finish()
	}()

	return s.next.Getpx(ctx, id)
}
