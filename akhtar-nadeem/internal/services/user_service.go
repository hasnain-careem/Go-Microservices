package services

import (
	"context"

	"github.com/akhtarCareem/golang-assignment/internal/repositories"
	"github.com/akhtarCareem/golang-assignment/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	proto.UnimplementedUserServiceServer
	store repositories.UserStore
	log   *logrus.Logger
}

func NewUserService(store repositories.UserStore, log *logrus.Logger) *UserServiceServer {
	return &UserServiceServer{store: store, log: log}
}

func (u *UserServiceServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	u.log.WithField("user_id", req.UserId).Info("GetUser called")
	name, err := u.store.GetUser(req.UserId)
	if err != nil {
		u.log.WithError(err).Error("failed to get user")
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &proto.GetUserResponse{Name: name}, nil
}

func (u *UserServiceServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	u.log.WithField("name", req.Name).Info("CreateUser called")
	userID, err := u.store.CreateUser(req.Name)
	if err != nil {
		u.log.WithError(err).Error("failed to create user")
		return nil, status.Error(codes.Internal, "failed to create user")
	}
	return &proto.CreateUserResponse{UserId: userID}, nil
}

func (u *UserServiceServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	u.log.WithField("user_id", req.UserId).Info("DeleteUser called")
	if err := u.store.DeleteUser(req.UserId); err != nil {
		u.log.WithError(err).Error("failed to delete user")
		return nil, status.Error(codes.Internal, "failed to delete user")
	}
	return &proto.DeleteUserResponse{Message: "user deleted successfully"}, nil
}
