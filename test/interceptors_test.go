package gerpc_test

import (
	"testing"

	"github.com/itsLeonB/gerpc"
)

func TestNewErrorInterceptor(t *testing.T) {
	logger := &mockLogger{}
	interceptor := gerpc.NewErrorInterceptor(logger)
	
	if interceptor == nil {
		t.Fatal("expected non-nil interceptor")
	}
}

func TestNewLoggingInterceptor(t *testing.T) {
	logger := &mockLogger{}
	interceptor := gerpc.NewLoggingInterceptor(logger)
	
	if interceptor == nil {
		t.Fatal("expected non-nil interceptor")
	}
}
