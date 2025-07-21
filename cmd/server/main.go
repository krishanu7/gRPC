package main

import (
	"github.com/krishanu7/grpc/internal/config"
	"github.com/krishanu7/grpc/internal/repository"
	"github.com/krishanu7/grpc/internal/server"
	"github.com/krishanu7/grpc/internal/service"
	"github.com/krishanu7/grpc/pkg/database"
	"github.com/krishanu7/grpc/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}
	logger, err := logger.NewLogger()

	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	db, err := database.NewDB(cfg)

	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, logger)
	grpcServer := server.NewGRPCServer(cfg, logger)

	if err := grpcServer.Start(userService); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
}