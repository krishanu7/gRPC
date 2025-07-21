package server

import (
	"net"

	"github.com/krishanu7/grpc/api/generated/user"
	"github.com/krishanu7/grpc/internal/config"
	"github.com/krishanu7/grpc/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	cfg    *config.Config
	logger *zap.Logger
}

func NewGRPCServer(cfg *config.Config, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *GRPCServer) Start(userService *service.UserService) error {
	lis, err := net.Listen("tcp", s.cfg.Port)
	if err != nil {
		s.logger.Fatal("Failed to listen", zap.Error(err))
		return err
	}
	grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(UnaryInterceptor(s.logger)),
        grpc.StreamInterceptor(StreamInterceptor(s.logger)),
    )
	user.RegisterUserServiceServer(grpcServer, userService)

	s.logger.Info("Starting gRPC server", zap.String("port", s.cfg.Port))
	return grpcServer.Serve(lis)
}
