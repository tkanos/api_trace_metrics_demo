//+build unit

package pixel

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockedService struct {
	mock.Mock
}

func (m *mockedService) Getpx(ctx context.Context, id string) (*px, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*px), args.Error(1)
}
