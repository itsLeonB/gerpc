package gerpc_test

import (
	"testing"

	"github.com/itsLeonB/gerpc"
	"github.com/stretchr/testify/assert"
)

func TestNewErrorInterceptor(t *testing.T) {
	logger := &MockLogger{}
	interceptor := gerpc.NewErrorInterceptor(logger)

	assert.NotNil(t, interceptor)
}

func TestNewLoggingInterceptor(t *testing.T) {
	logger := &MockLogger{}
	interceptor := gerpc.NewLoggingInterceptor(logger)

	assert.NotNil(t, interceptor)
}
