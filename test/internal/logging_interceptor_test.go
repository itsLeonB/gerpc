package internal_test

import (
	"context"
	"errors"
	"testing"

	"github.com/itsLeonB/gerpc/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLoggingInterceptor_Handle_Success(t *testing.T) {
	logger := &mockLogger{}
	interceptor := internal.NewLoggingInterceptor(logger)
	
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}
	
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	resp, err := interceptor.Handle(context.Background(), nil, info, handler)
	
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp != "success" {
		t.Errorf("expected 'success', got %v", resp)
	}
}

func TestLoggingInterceptor_Handle_Error(t *testing.T) {
	logger := &mockLogger{}
	interceptor := internal.NewLoggingInterceptor(logger)
	
	testErr := status.Error(codes.InvalidArgument, "test error")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, testErr
	}
	
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	_, err := interceptor.Handle(context.Background(), nil, info, handler)
	
	if err != testErr {
		t.Errorf("expected same error, got %v", err)
	}
}

func TestLoggingInterceptor_Handle_NonGRPCError(t *testing.T) {
	logger := &mockLogger{}
	interceptor := internal.NewLoggingInterceptor(logger)
	
	testErr := errors.New("regular error")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, testErr
	}
	
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	_, err := interceptor.Handle(context.Background(), nil, info, handler)
	
	if err != testErr {
		t.Errorf("expected same error, got %v", err)
	}
}
