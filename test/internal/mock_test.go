package internal_test

import (
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(args ...interface{}) { m.Called(args...) }
func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}
func (m *MockLogger) Info(args ...interface{}) { m.Called(args...) }
func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}
func (m *MockLogger) Warn(args ...interface{}) { m.Called(args...) }
func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}
func (m *MockLogger) Error(args ...interface{}) { m.Called(args...) }
func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}
func (m *MockLogger) Fatal(args ...interface{}) { m.Called(args...) }
func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

type mockServerStream struct{}

func (m *mockServerStream) SetHeader(metadata.MD) error  { return nil }
func (m *mockServerStream) SendHeader(metadata.MD) error { return nil }
func (m *mockServerStream) SetTrailer(metadata.MD)       {}
func (m *mockServerStream) Context() context.Context     { return context.Background() }
func (m *mockServerStream) SendMsg(interface{}) error    { return nil }
func (m *mockServerStream) RecvMsg(interface{}) error    { return nil }
