package internal_test

import (
	"context"
	"testing"

	"github.com/itsLeonB/gerpc/internal"
	"google.golang.org/grpc"
)

func TestInterceptorInterface(t *testing.T) {
	logger := &mockLogger{}
	
	// Test that error interceptor implements Interceptor interface
	var _ internal.Interceptor = internal.NewErrorInterceptor(logger)
	
	// Test that logging interceptor implements Interceptor interface
	var _ internal.Interceptor = internal.NewLoggingInterceptor(logger)
}

type testInterceptor struct{}

func (ti *testInterceptor) Handle(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	return handler(ctx, req)
}

func TestCustomInterceptor(t *testing.T) {
	var _ internal.Interceptor = &testInterceptor{}
}
