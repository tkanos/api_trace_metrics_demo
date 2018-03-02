package pixel

import (
	"context"
	"errors"
	"strconv"
	"time"
)

// ErrNotFound is used when an px is not found
var ErrNotFound = errors.New("px not found")

// Service is the Order service interface
type Service interface {
	Getpx(ctx context.Context, id string) (*Px, error)
}

type service struct {
	repository IRepository
}

// NewService return a new instance of order service
func NewService(repository IRepository) Service {
	return service{
		repository: repository,
	}
}

// Getpx returns an px regarding the id passed in parameter
func (s service) Getpx(ctx context.Context, id string) (m *Px, err error) {
	time.Sleep(50 * time.Millisecond)
	if _, errConv := strconv.Atoi(id); errConv != nil {
		s.repository.DoSomeWork(ctx, 700*time.Millisecond)
		return nil, ErrInvalidBody
	}

	s.repository.DoSomeWork(ctx, 100*time.Millisecond)

	m = &Px{PxID: id}
	return
}
