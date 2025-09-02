package gerpc

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsLeonB/ezutil/v2"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	logger          ezutil.Logger
	address         string
	opts            []grpc.ServerOption
	registerSrvFunc func(*grpc.Server) error
	shutdownFunc    func() error
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

func (s *GrpcServer) WithLogger(logger ezutil.Logger) *GrpcServer {
	s.logger = logger
	return s
}

func (s *GrpcServer) WithAddress(address string) *GrpcServer {
	s.address = address
	return s
}

func (s *GrpcServer) WithOpts(opts ...grpc.ServerOption) *GrpcServer {
	s.opts = opts
	return s
}

func (s *GrpcServer) WithRegisterSrvFunc(registerSrvFunc func(*grpc.Server) error) *GrpcServer {
	s.registerSrvFunc = registerSrvFunc
	return s
}

func (s *GrpcServer) WithShutdownFunc(shutdownFunc func() error) *GrpcServer {
	s.shutdownFunc = shutdownFunc
	return s
}

func (s *GrpcServer) Run() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		s.logger.Fatalf("error listening to %s: %v", s.address, err)
	}

	grpcServer := grpc.NewServer(s.opts...)
	if err := s.registerSrvFunc(grpcServer); err != nil {
		s.logger.Fatalf("error registering services: %v", err)
	}

	go func() {
		s.logger.Infof("server started at: %s", s.address)
		if err := grpcServer.Serve(listener); err != nil {
			s.logger.Fatalf("failed to serve: %v", err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	s.logger.Info("shutting down server...")
	grpcServer.GracefulStop()

	s.logger.Info("initating cleanup")
	if err := s.shutdownFunc(); err != nil {
		s.logger.Errorf("error during cleanup: %v", err)
	}

	s.logger.Info("server successfully shutdown")
}
