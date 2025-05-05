package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Booking struct {
	ID     int32
	UserID int32
	RideID int32
	Time   string
}

type BookingRepository interface {
	Create(ctx context.Context, userID, rideID int32) (*Booking, error)
	GetByID(ctx context.Context, id int32) (*Booking, error)
}

type PostgresBookingRepository struct {
	db *sql.DB
}

func NewPostgresBookingRepository(db *sql.DB) BookingRepository {
	return &PostgresBookingRepository{db: db}
}

func (r *PostgresBookingRepository) Create(ctx context.Context, userID, rideID int32) (*Booking, error) {
	var bookingID int32
	timestamp := time.Now().Format(time.RFC3339)
	query := `INSERT INTO bookings (user_id, ride_id, time) VALUES ($1, $2, $3) RETURNING booking_id`
	err := r.db.QueryRowContext(ctx, query, userID, rideID, timestamp).Scan(&bookingID)
	if err != nil {
		log.Printf("Create booking failed: %v", err)
		return nil, err
	}

	return &Booking{
		ID:     bookingID,
		UserID: userID,
		RideID: rideID,
		Time:   timestamp,
	}, nil
}

func (r *PostgresBookingRepository) GetByID(ctx context.Context, id int32) (*Booking, error) {
	query := `SELECT user_id, ride_id, time FROM bookings WHERE booking_id = $1`
	var userID, rideID int32
	var timeStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(&userID, &rideID, &timeStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking not found")
		}
		log.Printf("Get booking failed: %v", err)
		return nil, err
	}

	return &Booking{
		ID:     id,
		UserID: userID,
		RideID: rideID,
		Time:   timeStr,
	}, nil
}
