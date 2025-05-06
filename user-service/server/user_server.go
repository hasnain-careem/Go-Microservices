package server

import (
	"context"
	"fmt"

	pb "user-service/pb/proto/user"
	"user-service/repository"

	"github.com/hasnain-zafar/go-microservices/common/errors"
	"github.com/hasnain-zafar/go-microservices/common/logger"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	repo         repository.UserRepository
	logger       *logger.Logger
	errorHandler *errors.ErrorHandler
}

func NewUserServer(repo repository.UserRepository) *UserServer {
	log := logger.NewLogger("user-service")
	return &UserServer{
		repo:         repo,
		logger:       log,
		errorHandler: errors.NewErrorHandler(log),
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s.logger.LogRequest("CreateUser", req)

	if req.GetName() == "" {
		return nil, s.errorHandler.HandleInvalidArgument("invalid user name", fmt.Errorf("name cannot be empty"))
	}

	userID, err := s.repo.Create(ctx, req.GetName())
	if err != nil {
		return nil, s.errorHandler.HandleDatabaseError("failed to create user", err)
	}

	res := &pb.CreateUserResponse{UserId: userID}

	s.logger.LogResponse("CreateUser", res)

	return res, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	s.logger.LogRequest("GetUser", req)

	if req.GetUserId() <= 0 {
		return nil, s.errorHandler.HandleInvalidArgument("invalid user ID", fmt.Errorf("user ID must be positive"))
	}

	name, err := s.repo.GetByID(ctx, req.GetUserId())
	if err != nil {
		if err.Error() == "user not found" {
			return nil, s.errorHandler.HandleNotFound("user not found", err)
		}
		return nil, s.errorHandler.HandleDatabaseError("failed to get user", err)
	}

	res := &pb.GetUserResponse{Name: name}

	s.logger.LogResponse("GetUser", res)

	return res, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	s.logger.LogRequest("DeleteUser", req)

	if req.GetUserId() <= 0 {
		return nil, s.errorHandler.HandleInvalidArgument("invalid user ID", fmt.Errorf("user ID must be positive"))
	}

	message, err := s.repo.Delete(ctx, req.GetUserId())
	if err != nil {
		if err.Error() == "no user found to delete" {
			return nil, s.errorHandler.HandleNotFound("user not found", err)
		}
		return nil, s.errorHandler.HandleDatabaseError("failed to delete user", err)
	}

	res := &pb.DeleteUserResponse{Message: message}

	s.logger.LogResponse("DeleteUser", res)

	return res, nil
}
