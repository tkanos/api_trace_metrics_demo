package pixel

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints represent the order service endpoints
type Endpoints struct {
	GetByID endpoint.Endpoint
}

// MakeGetByIDEndpoint returns an endpoint used for getting one px
func MakeGetByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetpxRequest)

		return s.Getpx(ctx, req.ID)
	}
}

// GetpxRequest represents the request parameters used for getting one px
type GetpxRequest struct {
	ID string `json:"id"`
}
