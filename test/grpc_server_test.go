package gerpc_test

import (
	"testing"

	"github.com/itsLeonB/gerpc"
	"google.golang.org/grpc"
)

func TestNewGrpcServer(t *testing.T) {
	server := gerpc.NewGrpcServer()
	if server == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestGrpcServer_WithLogger(t *testing.T) {
	server := gerpc.NewGrpcServer()
	logger := &mockLogger{}
	
	result := server.WithLogger(logger)
	if result != server {
		t.Error("expected fluent interface")
	}
}

func TestGrpcServer_WithAddress(t *testing.T) {
	server := gerpc.NewGrpcServer()
	address := ":8080"
	
	result := server.WithAddress(address)
	if result != server {
		t.Error("expected fluent interface")
	}
}

func TestGrpcServer_WithOpts(t *testing.T) {
	server := gerpc.NewGrpcServer()
	opts := []grpc.ServerOption{}
	
	result := server.WithOpts(opts...)
	if result != server {
		t.Error("expected fluent interface")
	}
}

func TestGrpcServer_WithRegisterSrvFunc(t *testing.T) {
	server := gerpc.NewGrpcServer()
	registerFunc := func(*grpc.Server) error { return nil }
	
	result := server.WithRegisterSrvFunc(registerFunc)
	if result != server {
		t.Error("expected fluent interface")
	}
}

func TestGrpcServer_WithShutdownFunc(t *testing.T) {
	server := gerpc.NewGrpcServer()
	shutdownFunc := func() error { return nil }
	
	result := server.WithShutdownFunc(shutdownFunc)
	if result != server {
		t.Error("expected fluent interface")
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
