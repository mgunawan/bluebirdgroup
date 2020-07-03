package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

//AddMetadata ...
func AddMetadata(ctx context.Context, key, value string) context.Context {
	md := metadata.Pairs(key, value)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}
