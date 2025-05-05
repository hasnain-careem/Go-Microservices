package server

import (
	"booking-service/pb/proto/booking"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	userpb "user-service/pb/proto/user"
	ridepb "ride-service/pb/proto/ride"
)

type BookingServer struct {
	pb.UnimplementedBookingServiceServer
	DB         *sql.DB
	UserClient userpb.UserServiceClient
	RideClient ridepb.RideServiceClient
}

func (s *BookingServer) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.Booking, error) {
	// Insert ride by calling RideService (assumes CreateRide exists)
	rideReq := &ridepb.CreateRideRequest{
		Source:      req.Ride.Source,
		Destination: req.Ride.Destination,
		Distance:    req.Ride.Distance,
		Cost:        req.Ride.Cost,
	}
	rideRes, err := s.RideClient.CreateRide(ctx, rideReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create ride: %v", err)
	}

	// Insert booking in local DB
	var bookingID int32
	timestamp := time.Now().Format(time.RFC3339)
	query := `INSERT INTO bookings (user_id, ride_id, time) VALUES ($1, $2, $3) RETURNING booking_id`
	err = s.DB.QueryRowContext(ctx, query, req.UserId, rideRes.RideId, timestamp).Scan(&bookingID)
	if err != nil {
		log.Printf("failed to insert booking: %v", err)
		return nil, err
	}

	return &pb.Booking{
		BookingId: bookingID,
		UserId:    req.UserId,
		RideId:    rideRes.RideId,
		Time:      timestamp,
	}, nil
}

func (s *BookingServer) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.BookingDetails, error) {
	// Step 1: Get booking from local DB
	query := `SELECT user_id, ride_id, time FROM bookings WHERE booking_id = $1`
	var userID, rideID int32
	var timeStr string

	err := s.DB.QueryRowContext(ctx, query, req.BookingId).Scan(&userID, &rideID, &timeStr)
	if err != nil {
		return nil, fmt.Errorf("booking not found: %v", err)
	}

	// Step 2: Get user name from UserService
	userRes, err := s.UserClient.GetUser(ctx, &userpb.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Step 3: Get ride details from RideService
	rideRes, err := s.RideClient.GetRide(ctx, &ridepb.GetRideRequest{RideId: rideID})
	if err != nil {
		return nil, fmt.Errorf("failed to get ride: %v", err)
	}

	return &pb.BookingDetails{
		Name:        userRes.Name,
		Source:      rideRes.Source,
		Destination: rideRes.Destination,
		Distance:    rideRes.Distance,
		Cost:        rideRes.Cost,
		Time:        timeStr,
	}, nil
}
