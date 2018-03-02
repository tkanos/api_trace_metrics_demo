package pixel

import (
	"context"
	"time"

	tracing "github.com/ricardo-ch/go-tracing"
)

//IRepository represent a fake third party repository
type IRepository interface {
	DoSomeWork(ctx context.Context, d time.Duration) error
}

type repository struct{}

// NewRepository ...
func NewRepository() IRepository {
	return repository{}
}

func (r repository) DoSomeWork(ctx context.Context, d time.Duration) error {
	time.Sleep(d)
	return nil
}

type tracingRepository struct {
	next IRepository
}

// NewTracingRepository ...
func NewTracingRepository(repository IRepository) IRepository {
	return tracingRepository{
		next: repository,
	}
}

func (r tracingRepository) DoSomeWork(ctx context.Context, d time.Duration) (err error) {
	span, ctx := tracing.CreateSpan(ctx, "pxs.Repository::DoSomeWork", &map[string]interface{}{"duration": d})
	defer func() {
		if err != nil {
			tracing.SetSpanError(span, err)
		}
		span.Finish()
	}()

	return r.next.DoSomeWork(ctx, d)
}
