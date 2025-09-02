package internal_test

import (
	"context"
	"errors"
	"testing"

	"github.com/itsLeonB/gerpc/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLoggingInterceptor_Handle_Success(t *testing.T) {
	logger := &MockLogger{}
	logger.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return()

	interceptor := internal.NewLoggingInterceptor(logger)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	resp, err := interceptor.Handle(context.Background(), nil, info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
	logger.AssertExpectations(t)
}

func TestLoggingInterceptor_Handle_Error(t *testing.T) {
	logger := &MockLogger{}
	logger.On("Errorf", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	interceptor := internal.NewLoggingInterceptor(logger)

	testErr := status.Error(codes.InvalidArgument, "test error")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, testErr
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	_, err := interceptor.Handle(context.Background(), nil, info, handler)

	assert.Equal(t, testErr, err)
	logger.AssertExpectations(t)
}

func TestLoggingInterceptor_Handle_NonGRPCError(t *testing.T) {
	logger := &MockLogger{}
	logger.On("Errorf", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	interceptor := internal.NewLoggingInterceptor(logger)

	testErr := errors.New("regular error")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, testErr
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	_, err := interceptor.Handle(context.Background(), nil, info, handler)

	assert.Equal(t, testErr, err)
	logger.AssertExpectations(t)
}
