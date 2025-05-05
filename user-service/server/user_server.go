package server

import (
	"context"
	"log"

	pb "user-service/pb/proto/user"
	"user-service/repository"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	repo repository.UserRepository
}

func NewUserServer(repo repository.UserRepository) *UserServer {
	return &UserServer{repo: repo}
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userID, err := s.repo.Create(ctx, req.GetName())
	if err != nil {
		log.Printf("CreateUser failed: %v", err)
		return nil, err
	}
	return &pb.CreateUserResponse{UserId: userID}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	name, err := s.repo.GetByID(ctx, req.GetUserId())
	if err != nil {
		log.Printf("GetUser failed: %v", err)
		return nil, err
	}
	return &pb.GetUserResponse{Name: name}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	message, err := s.repo.Delete(ctx, req.GetUserId())
	if err != nil {
		log.Printf("DeleteUser failed: %v", err)
		return nil, err
	}

	return &pb.DeleteUserResponse{Message: message}, nil
}
