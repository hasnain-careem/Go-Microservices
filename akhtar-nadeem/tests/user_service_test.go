package tests

import (
	"context"
	"github.com/akhtarCareem/golang-assignment/internal/services"
	proto "github.com/akhtarCareem/golang-assignment/proto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockUserStore struct {
	users map[string]string
	err   error
}

func (m *mockUserStore) CreateUser(name string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	id := "mock_user_id"
	m.users[id] = name
	return id, nil
}
func (m *mockUserStore) GetUser(userID string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.users[userID], nil
}
func (m *mockUserStore) DeleteUser(userID string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.users, userID)
	return nil
}

func TestUserService(t *testing.T) {
	store := &mockUserStore{users: make(map[string]string)}
	log := logrus.New()
	svc := services.NewUserService(store, log)

	// CreateUser
	respCreate, err := svc.CreateUser(context.Background(), &proto.CreateUserRequest{Name: "John"})
	assert.NoError(t, err)
	assert.NotEmpty(t, respCreate.UserId)

	// GetUser
	respGet, err := svc.GetUser(context.Background(), &proto.GetUserRequest{UserId: respCreate.UserId})
	assert.NoError(t, err)
	assert.Equal(t, "John", respGet.Name)

	// DeleteUser
	respDel, err := svc.DeleteUser(context.Background(), &proto.DeleteUserRequest{UserId: respCreate.UserId})
	assert.NoError(t, err)
	assert.Equal(t, "user deleted successfully", respDel.Message)
}
