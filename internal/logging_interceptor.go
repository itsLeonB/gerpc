package internal

import (
	"context"
	"time"

	"github.com/itsLeonB/ezutil/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type loggingInterceptor struct {
	logger ezutil.Logger
}

func NewLoggingInterceptor(logger ezutil.Logger) Interceptor {
	return &loggingInterceptor{logger}
}

func (li *loggingInterceptor) Handle(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()

	// Call handler
	resp, err = handler(ctx, req)

	// Duration
	elapsed := time.Since(start)

	// Extract gRPC status code (if error)
	st, _ := status.FromError(err)

	if err != nil {
		li.logger.Errorf(
			"[gRPC] method=%s duration=%s status=%s error=%v",
			info.FullMethod,
			elapsed,
			st.Code(),
			st.Message(),
		)
	} else {
		li.logger.Infof(
			"[gRPC] method=%s duration=%s status=OK",
			info.FullMethod,
			elapsed,
		)
	}

	return resp, err
}
