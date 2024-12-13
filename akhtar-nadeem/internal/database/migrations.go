package database

import "github.com/jinzhu/gorm"

// User model
type User struct {
	UserID string `gorm:"primary_key"`
	Name   string
}

// Ride model
type Ride struct {
	RideID      string `gorm:"primary_key"`
	Source      string
	Destination string
	Distance    int32
	Cost        int32
}

// Booking model
type Booking struct {
	BookingID string `gorm:"primary_key"`
	UserID    string
	RideID    string
	Time      string
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Ride{}, &Booking{}).Error
	if err != nil {
		return err
	}
	return nil
}
