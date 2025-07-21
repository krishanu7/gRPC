package server

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Logging
		logger.Info("Unary call", zap.String("method", info.FullMethod))

		// Authentication (commented out)
		// if err := authenticate(ctx); err != nil {
		//     logger.Error("Authentication failed", zap.Error(err))
		//     return nil, err
		// }

		// Error Handling
		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error("Unary call failed",
				zap.String("method", info.FullMethod),
				zap.Error(err))
			return nil, wrapError(err)
		}
		return resp, nil
	}
}

func StreamInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Logging
		logger.Info("Stream call", zap.String("method", info.FullMethod))

		// Authentication (commented out)
		// if err := authenticate(stream.Context()); err != nil {
		//     logger.Error("Authentication failed", zap.Error(err))
		//     return err
		// }

		// Error Handling
		err := handler(srv, stream)
		if err != nil {
			logger.Error("Stream call failed",
				zap.String("method", info.FullMethod),
				zap.Error(err))
			return wrapError(err)
		}
		return nil
	}
}

func authenticate(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}
	token, ok := md["authorizatioin"]

	if !ok || len(token) == 0 {
		return status.Error(codes.Unauthenticated, "missing ")
	}
	if !strings.HasPrefix(token[0], "Bearer ") {
		return status.Error(codes.Unauthenticated, "invalid token format")
	}
	if token[0] != "Bearer valid-token" {
		return status.Error(codes.Unauthenticated, "invalid token")
	}
	return nil
}

func wrapError(err error) error {
	if _, ok := status.FromError(err); ok {
		return err
	}
	return status.Errorf(codes.Internal, "internal error: %v", err)
}
