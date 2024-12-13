package repositories

import (
	"github.com/akhtarCareem/golang-assignment/internal/database"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserStore interface {
	CreateUser(name string) (string, error)
	GetUser(userID string) (string, error)
	DeleteUser(userID string) error
}

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{db: db}
}

func (u *userStore) CreateUser(name string) (string, error) {
	id := uuid.New().String()
	user := database.User{UserID: id, Name: name}
	if err := u.db.Create(&user).Error; err != nil {
		return "", err
	}
	return id, nil
}

func (u *userStore) GetUser(userID string) (string, error) {
	var user database.User
	if err := u.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Name, nil
}

func (u *userStore) DeleteUser(userID string) error {
	return u.db.Where("user_id = ?", userID).Delete(&database.User{}).Error
}
