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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorInterceptor_Handle_Success(t *testing.T) {
	logger := &mockLogger{}
	interceptor := internal.NewErrorInterceptor(logger)
	
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}
	
	resp, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp != "success" {
		t.Errorf("expected 'success', got %v", resp)
	}
}

func TestErrorInterceptor_Handle_AppError(t *testing.T) {
	logger := &mockLogger{}
	interceptor := internal.NewErrorInterceptor(logger)
	
	appErr := ungerr.BadRequestError("test error")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, appErr
	}
	
	_, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status error")
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}
}

func TestErrorInterceptor_Handle_ValidationError(t *testing.T) {
	logger := &mockLogger{}
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
	if !ok {
		t.Fatal("expected gRPC status error")
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}
}

func TestErrorInterceptor_Handle_JSONError(t *testing.T) {
	logger := &mockLogger{}
	interceptor := internal.NewErrorInterceptor(logger)
	
	jsonErr := &json.SyntaxError{}
	wrappedErr := eris.Wrap(jsonErr, "json error")
	
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, wrappedErr
	}
	
	_, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status error")
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}
}

func TestErrorInterceptor_Handle_EOFError(t *testing.T) {
	logger := &mockLogger{}
	interceptor := internal.NewErrorInterceptor(logger)
	
	wrappedErr := eris.Wrap(io.EOF, "eof error")
	
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, wrappedErr
	}
	
	_, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status error")
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}
}

func TestErrorInterceptor_Handle_Panic(t *testing.T) {
	logger := &mockLogger{}
	interceptor := internal.NewErrorInterceptor(logger)
	
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test panic")
	}
	
	_, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status error")
	}
	if st.Code() != codes.Internal {
		t.Errorf("expected Internal, got %v", st.Code())
	}
}

func TestErrorInterceptor_Handle_UnwrappedError(t *testing.T) {
	logger := &mockLogger{}
	interceptor := internal.NewErrorInterceptor(logger)
	
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("unwrapped error")
	}
	
	_, err := interceptor.Handle(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status error")
	}
	if st.Code() != codes.Internal {
		t.Errorf("expected Internal, got %v", st.Code())
	}
}

type mockLogger struct{}

func (m *mockLogger) Debug(args ...interface{})                 {}
func (m *mockLogger) Debugf(format string, args ...interface{}) {}
func (m *mockLogger) Info(args ...interface{})                  {}
func (m *mockLogger) Infof(format string, args ...interface{})  {}
func (m *mockLogger) Warn(args ...interface{})                  {}
func (m *mockLogger) Warnf(format string, args ...interface{})  {}
func (m *mockLogger) Error(args ...interface{})                 {}
func (m *mockLogger) Errorf(format string, args ...interface{}) {}
func (m *mockLogger) Fatal(args ...interface{})                 {}
func (m *mockLogger) Fatalf(format string, args ...interface{}) {}
