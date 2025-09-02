package internal_test

import (
	"context"
	"testing"

	"github.com/itsLeonB/gerpc/internal"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestInterceptorInterface(t *testing.T) {
	logger := &MockLogger{}

	errorInterceptor := internal.NewErrorInterceptor(logger)
	loggingInterceptor := internal.NewLoggingInterceptor(logger)

	assert.NotNil(t, errorInterceptor)
	assert.NotNil(t, loggingInterceptor)
}

type testInterceptor struct{}

func (ti *testInterceptor) Handle(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	return handler(ctx, req)
}

func TestCustomInterceptor(t *testing.T) {
	var interceptor internal.Interceptor = &testInterceptor{}
	assert.NotNil(t, interceptor)
}
