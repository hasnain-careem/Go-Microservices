package repositories

import (
	"time"

	"github.com/akhtarCareem/golang-assignment/internal/database"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type BookingRepository interface {
	CreateBooking(userID string, source, destination string, distance, cost int32) (string, string, time.Time, error)
	GetBookingDetails(bookingID string) (name, source, destination string, distance, cost int32, bookingTime time.Time, userID, rideID string, err error)
}

type bookingRepo struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepo{db: db}
}

func (b *bookingRepo) CreateBooking(userID, source, destination string, distance, cost int32) (string, string, time.Time, error) {
	rideID := uuid.New().String()
	ride := database.Ride{
		RideID:      rideID,
		Source:      source,
		Destination: destination,
		Distance:    distance,
		Cost:        cost,
	}
	if err := b.db.Create(&ride).Error; err != nil {
		return "", "", time.Time{}, err
	}
	bookingID := uuid.New().String()
	now := time.Now()
	booking := database.Booking{
		BookingID: bookingID,
		UserID:    userID,
		RideID:    rideID,
		Time:      now.Format(time.RFC3339),
	}
	if err := b.db.Create(&booking).Error; err != nil {
		return "", "", time.Time{}, err
	}

	return bookingID, rideID, now, nil
}

func (b *bookingRepo) GetBookingDetails(bookingID string) (name, source, destination string, distance, cost int32, bookingTime time.Time, userID, rideID string, err error) {
	var booking database.Booking
	if err = b.db.Where("booking_id = ?", bookingID).First(&booking).Error; err != nil {
		return
	}

	var user database.User
	if err = b.db.Where("user_id = ?", booking.UserID).First(&user).Error; err != nil {
		return
	}

	var ride database.Ride
	if err = b.db.Where("ride_id = ?", booking.RideID).First(&ride).Error; err != nil {
		return
	}

	t, err := time.Parse(time.RFC3339, booking.Time)
	if err != nil {
		return
	}

	return user.Name, ride.Source, ride.Destination, ride.Distance, ride.Cost, t, booking.UserID, booking.RideID, nil
}
