package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

//RequestValidator ...
type RequestValidator interface {
	Validate() error
}

//ValidateMiddleware validate grpc request based on https://github.com/mwitkow/go-proto-validators
func ValidateMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			req := request.(RequestValidator)
			if err := req.Validate(); err != nil {
				return nil, err
			}
			return next(ctx, request)
		}
	}
}
