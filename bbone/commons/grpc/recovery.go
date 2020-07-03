package grpc

import (
	"errors"

	"bluebirdgroup/bbone/commons/logger"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

//Recovery return grpc server option with recovery handler
func Recovery() []grpc.ServerOption {
	handler := func(p interface{}) (err error) {
		logger.Error(p)
		return errors.New("Internal Server Error")
	}
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(handler),
	}
	serverOptions := []grpc.ServerOption{
		middleware.WithUnaryServerChain(
			recovery.UnaryServerInterceptor(opts...),
		),
		middleware.WithStreamServerChain(
			recovery.StreamServerInterceptor(opts...),
		)}
	return serverOptions
}
