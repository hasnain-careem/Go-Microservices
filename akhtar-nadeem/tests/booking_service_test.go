package tests

import (
	"context"
	"errors"
	"github.com/akhtarCareem/golang-assignment/internal/services"
	proto "github.com/akhtarCareem/golang-assignment/proto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockBookingRepo struct {
	// For simplicity, we’ll store bookings and rides in maps keyed by ID.
	bookings map[string]mockBooking
	rides    map[string]mockRide
	err      error
}

type mockBooking struct {
	BookingID string
	UserID    string
	RideID    string
	Time      time.Time
}

type mockRide struct {
	RideID      string
	Source      string
	Destination string
	Distance    int32
	Cost        int32
}

func (m *mockBookingRepo) CreateBooking(userID, source, destination string, distance, cost int32) (string, string, time.Time, error) {
	if m.err != nil {
		return "", "", time.Time{}, m.err
	}
	rideID := "mock_ride_id"
	bookingID := "mock_booking_id"
	now := time.Now()

	m.rides[rideID] = mockRide{RideID: rideID, Source: source, Destination: destination, Distance: distance, Cost: cost}
	m.bookings[bookingID] = mockBooking{BookingID: bookingID, UserID: userID, RideID: rideID, Time: now}
	return bookingID, rideID, now, nil
}

func (m *mockBookingRepo) GetBookingDetails(bookingID string) (name, source, destination string, distance, cost int32, bookingTime time.Time, userID, rideID string, err error) {
	if m.err != nil {
		return "", "", "", 0, 0, time.Time{}, "", "", m.err
	}
	b, ok := m.bookings[bookingID]
	if !ok {
		return "", "", "", 0, 0, time.Time{}, "", "", errors.New("not found")
	}
	r, ok := m.rides[b.RideID]
	if !ok {
		return "", "", "", 0, 0, time.Time{}, "", "", errors.New("ride not found")
	}
	// For simplicity, let's assume the user name is not stored here; just return a hard-coded name "John"
	return "John", r.Source, r.Destination, r.Distance, r.Cost, b.Time, b.UserID, b.RideID, nil
}

func TestBookingService(t *testing.T) {
	repo := &mockBookingRepo{
		bookings: make(map[string]mockBooking),
		rides:    make(map[string]mockRide),
	}
	log := logrus.New()
	svc := services.NewBookingService(repo, log)

	// CreateBooking
	respCreate, err := svc.CreateBooking(context.Background(), &proto.CreateBookingRequest{
		UserId: "mock_user_id",
		Ride: &proto.Ride{
			Source:      "A",
			Destination: "B",
			Distance:    10,
			Cost:        100,
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, respCreate.BookingId)
	assert.NotEmpty(t, respCreate.RideId)
	assert.NotEmpty(t, respCreate.UserId)
	assert.NotNil(t, respCreate.Time)

	// GetBooking
	respGet, err := svc.GetBooking(context.Background(), &proto.GetBookingRequest{BookingId: respCreate.BookingId})
	assert.NoError(t, err)
	assert.Equal(t, "John", respGet.Name)
	assert.Equal(t, "A", respGet.Source)
	assert.Equal(t, "B", respGet.Destination)
	assert.Equal(t, int32(10), respGet.Distance)
	assert.Equal(t, int32(100), respGet.Cost)
	assert.Equal(t, "mock_user_id", respGet.UserId)
	assert.Equal(t, respCreate.BookingId, respGet.BookingId)
	assert.Equal(t, respCreate.RideId, respGet.RideId)
	assert.NotNil(t, respGet.Time)

	// Error scenario
	repo.err = errors.New("db error")
	respGet, err = svc.GetBooking(context.Background(), &proto.GetBookingRequest{BookingId: "non_existent"})
	assert.Error(t, err)
	assert.Nil(t, respGet)
}
