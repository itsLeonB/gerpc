package internal_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/gerpc/internal"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorInterceptor_Handle_Success(t *testing.T) {
	logger := &MockLogger{}
	interceptor := internal.NewErrorInterceptor(logger)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	resp, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
}

func TestErrorInterceptor_Handle_AppError(t *testing.T) {
	logger := &MockLogger{}
	interceptor := internal.NewErrorInterceptor(logger)

	appErr := ungerr.BadRequestError("test error")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, appErr
	}

	_, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestErrorInterceptor_Handle_ValidationError(t *testing.T) {
	logger := &MockLogger{}
	interceptor := internal.NewErrorInterceptor(logger)

	validate := validator.New()
	type testStruct struct {
		Email string `validate:"required,email"`
	}

	validationErr := validate.Struct(&testStruct{Email: "invalid"})
	wrappedErr := eris.Wrap(validationErr, "validation failed")

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, wrappedErr
	}

	_, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestErrorInterceptor_Handle_JSONError(t *testing.T) {
	logger := &MockLogger{}
	interceptor := internal.NewErrorInterceptor(logger)

	jsonErr := &json.SyntaxError{}
	wrappedErr := eris.Wrap(jsonErr, "json error")

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, wrappedErr
	}

	_, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestErrorInterceptor_Handle_EOFError(t *testing.T) {
	logger := &MockLogger{}
	interceptor := internal.NewErrorInterceptor(logger)

	wrappedErr := eris.Wrap(io.EOF, "eof error")

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, wrappedErr
	}

	_, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestErrorInterceptor_Handle_Panic(t *testing.T) {
	logger := &MockLogger{}
	logger.On("Error", mock.Anything).Return()
	logger.On("Errorf", mock.Anything, mock.Anything).Return()

	interceptor := internal.NewErrorInterceptor(logger)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test panic")
	}

	_, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestErrorInterceptor_Handle_UnwrappedError(t *testing.T) {
	logger := &MockLogger{}
	logger.On("Error", mock.Anything).Return()
	logger.On("Errorf", mock.Anything, mock.Anything).Return()

	interceptor := internal.NewErrorInterceptor(logger)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("unwrapped error")
	}

	_, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}
