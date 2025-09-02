package gerpc_test

import (
	"testing"

	"github.com/itsLeonB/gerpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestNewGrpcServer(t *testing.T) {
	server := gerpc.NewGrpcServer()
	assert.NotNil(t, server)
}

func TestGrpcServer_WithLogger(t *testing.T) {
	server := gerpc.NewGrpcServer()
	logger := &MockLogger{}

	result := server.WithLogger(logger)
	assert.Equal(t, server, result)
}

func TestGrpcServer_WithAddress(t *testing.T) {
	server := gerpc.NewGrpcServer()

	result := server.WithAddress(":8080")
	assert.Equal(t, server, result)
}

func TestGrpcServer_WithOpts(t *testing.T) {
	server := gerpc.NewGrpcServer()

	result := server.WithOpts()
	assert.Equal(t, server, result)
}

func TestGrpcServer_WithRegisterSrvFunc(t *testing.T) {
	server := gerpc.NewGrpcServer()

	result := server.WithRegisterSrvFunc(func(*grpc.Server) error { return nil })
	assert.Equal(t, server, result)
}

func TestGrpcServer_WithShutdownFunc(t *testing.T) {
	server := gerpc.NewGrpcServer()

	result := server.WithShutdownFunc(func() error { return nil })
	assert.Equal(t, server, result)
}
