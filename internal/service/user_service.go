package service

import (
	"context"
	"io"

	"github.com/krishanu7/grpc/api/generated/user"
	"github.com/krishanu7/grpc/internal/model"
	"github.com/krishanu7/grpc/internal/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	user.UnimplementedUserServiceServer
	repo *repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo *repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func (s *UserService) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	s.logger.Info("Creating user", zap.String("email", req.Email))
	u := &model.User{Name: req.Name, Email: req.Email}
	id, err := s.repo.Create(u);
	
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
        return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	return &user.CreateUserResponse{
		Id: id,
		Name: req.Name,
		Email: req.Email,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	s.logger.Info("Getting user", zap.String("id", req.Id))
	u, err := s.repo.Get(req.Id)
	if err != nil {
		s.logger.Error("Failed to get user", zap.Error(err))
        return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}
	return &user.GetUserResponse{
		Id: u.ID,
		Name: u.Name,
		Email: u.Email,
	}, nil
}

func (s *UserService) UploadUsers(stream user.UserService_UploadUsersServer) error {
	s.logger.Info("Starting client-streaming uploadusers")
	var ids []string
	count := 0
	for {
		req, err := stream.Recv()
		
		if err == io.EOF {
			s.logger.Info("Completed client-streaming uploadUsers", zap.Int("count", count))
			return stream.SendAndClose(&user.UploadUsersResponse{
				Count: int32(count),
				Ids: ids,
			})
		}
		if err != nil {
			s.logger.Error("Error receiving stream", zap.Error(err))
            return status.Errorf(codes.InvalidArgument, "stream error: %v", err)
		}

		u := &model.User{Name: req.Name, Email: req.Email}
		id, err := s.repo.Create(u)

		if err != nil {
            s.logger.Error("Failed to create user in stream", zap.Error(err))
            return status.Errorf(codes.Internal, "failed to create user: %v", err)
        }

		ids = append(ids, id)
		count++
	}
}