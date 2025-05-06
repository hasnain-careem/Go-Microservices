package server

import (
	"context"
	"errors"
	"testing"

	pb "user-service/pb/proto/user"
	"user-service/repository/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	// Setup
	mockRepo := new(mocks.UserRepository)
	userServer := NewUserServer(mockRepo)

	ctx := context.Background()
	req := &pb.CreateUserRequest{Name: "John Doe"}

	// Expectations
	mockRepo.On("Create", ctx, "John Doe").Return(int32(1), nil)

	// Action
	resp, err := userServer.CreateUser(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.UserId)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_EmptyName(t *testing.T) {
	// Setup
	mockRepo := new(mocks.UserRepository)
	userServer := NewUserServer(mockRepo)

	ctx := context.Background()
	req := &pb.CreateUserRequest{Name: ""}

	// Action
	resp, err := userServer.CreateUser(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	// Repository should not be called when validation fails
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateUser_RepositoryError(t *testing.T) {
	// Setup
	mockRepo := new(mocks.UserRepository)
	userServer := NewUserServer(mockRepo)

	ctx := context.Background()
	req := &pb.CreateUserRequest{Name: "John Doe"}

	// Expectations
	mockRepo.On("Create", ctx, "John Doe").Return(int32(0), errors.New("database error"))

	// Action
	resp, err := userServer.CreateUser(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	// Setup
	mockRepo := new(mocks.UserRepository)
	userServer := NewUserServer(mockRepo)

	ctx := context.Background()
	req := &pb.GetUserRequest{UserId: 1}

	// Expectations
	mockRepo.On("GetByID", ctx, int32(1)).Return("John Doe", nil)

	// Action
	resp, err := userServer.GetUser(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "John Doe", resp.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_InvalidId(t *testing.T) {
	// Setup
	mockRepo := new(mocks.UserRepository)
	userServer := NewUserServer(mockRepo)

	ctx := context.Background()
	req := &pb.GetUserRequest{UserId: 0}

	// Action
	resp, err := userServer.GetUser(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	// Repository should not be called when validation fails
	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestGetUser_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(mocks.UserRepository)
	userServer := NewUserServer(mockRepo)

	ctx := context.Background()
	req := &pb.GetUserRequest{UserId: 999}

	// Expectations
	mockRepo.On("GetByID", ctx, int32(999)).Return("", errors.New("user not found"))

	// Action
	resp, err := userServer.GetUser(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	// Setup
	mockRepo := new(mocks.UserRepository)
	userServer := NewUserServer(mockRepo)

	ctx := context.Background()
	req := &pb.DeleteUserRequest{UserId: 1}
	expectedMsg := "User with ID 1 deleted successfully"

	// Expectations
	mockRepo.On("Delete", ctx, int32(1)).Return(expectedMsg, nil)

	// Action
	resp, err := userServer.DeleteUser(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedMsg, resp.Message)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_InvalidId(t *testing.T) {
	// Setup
	mockRepo := new(mocks.UserRepository)
	userServer := NewUserServer(mockRepo)

	ctx := context.Background()
	req := &pb.DeleteUserRequest{UserId: 0}

	// Action
	resp, err := userServer.DeleteUser(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	// Repository should not be called when validation fails
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestDeleteUser_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(mocks.UserRepository)
	userServer := NewUserServer(mockRepo)

	ctx := context.Background()
	req := &pb.DeleteUserRequest{UserId: 999}

	// Expectations
	mockRepo.On("Delete", ctx, int32(999)).Return("", errors.New("no user found to delete"))

	// Action
	resp, err := userServer.DeleteUser(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}
