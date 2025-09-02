package internal

import (
	"context"

	"google.golang.org/grpc"
)

type Interceptor interface {
	Handle(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error)
}
