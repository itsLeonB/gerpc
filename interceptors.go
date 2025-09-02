package gerpc

import (
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/gerpc/internal"
	"google.golang.org/grpc"
)

// NewErrorInterceptor creates an error handling interceptor for gRPC.
// It captures errors and panics from gRPC handlers, converts them into
// appropriate gRPC status codes with structured error messages.
func NewErrorInterceptor(logger ezutil.Logger) grpc.UnaryServerInterceptor {
	interceptor := internal.NewErrorInterceptor(logger)
	return interceptor.Handle
}

// NewLoggingInterceptor logs incoming requests, responses, durations, and errors.
func NewLoggingInterceptor(logger ezutil.Logger) grpc.UnaryServerInterceptor {
	interceptor := internal.NewLoggingInterceptor(logger)
	return interceptor.Handle
}
