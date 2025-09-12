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
			"[gRPC] method=%s duration=%v status=%s msg=%q err=%v",
			info.FullMethod,
			elapsed,
			st.Code().String(),
			st.Message(),
			err,
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

func (li *loggingInterceptor) HandleStream(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()

	err := handler(srv, ss)

	elapsed := time.Since(start)
	st, _ := status.FromError(err)

	if err != nil {
		li.logger.Errorf(
			"[gRPC Stream] method=%s duration=%v status=%s msg=%q err=%v",
			info.FullMethod,
			elapsed,
			st.Code().String(),
			st.Message(),
			err,
		)
	} else {
		li.logger.Infof(
			"[gRPC Stream] method=%s duration=%s status=OK",
			info.FullMethod,
			elapsed,
		)
	}

	return err
}
