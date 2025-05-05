package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	pb "user-service/pb/proto/user"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	DB *sql.DB
}

// CreateUser inserts a new user and returns the new user_id
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	query := `INSERT INTO users (name) VALUES ($1) RETURNING user_id`
	var userID int32
	err := s.DB.QueryRowContext(ctx, query, req.GetName()).Scan(&userID)
	if err != nil {
		log.Printf("CreateUser failed: %v", err)
		return nil, err
	}
	return &pb.CreateUserResponse{UserId: userID}, nil
}

// GetUser fetches a user by ID and returns the name
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	query := `SELECT name FROM users WHERE user_id = $1`
	var name string
	err := s.DB.QueryRowContext(ctx, query, req.GetUserId()).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("GetUser failed: %v", err)
		return nil, err
	}
	return &pb.GetUserResponse{Name: name}, nil
}

// DeleteUser removes a user by ID
func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	query := `DELETE FROM users WHERE user_id = $1`
	res, err := s.DB.ExecContext(ctx, query, req.GetUserId())
	if err != nil {
		log.Printf("DeleteUser failed: %v", err)
		return nil, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no user found to delete")
	}

	msg := fmt.Sprintf("User with ID %d deleted successfully", req.GetUserId())
	return &pb.DeleteUserResponse{Message: msg}, nil
}
