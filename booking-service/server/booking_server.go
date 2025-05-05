package server

import (
	pb "booking-service/pb/proto/booking"
	"booking-service/repository"
	"context"
	"fmt"
	"log"

	ridepb "ride-service/pb/proto/ride"
	userpb "user-service/pb/proto/user"
)

type BookingServer struct {
	pb.UnimplementedBookingServiceServer
	repo       repository.BookingRepository
	userClient userpb.UserServiceClient
	rideClient ridepb.RideServiceClient
}

func NewBookingServer(
	repo repository.BookingRepository,
	userClient userpb.UserServiceClient,
	rideClient ridepb.RideServiceClient,
) *BookingServer {
	return &BookingServer{
		repo:       repo,
		userClient: userClient,
		rideClient: rideClient,
	}
}

func (s *BookingServer) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.Booking, error) {
	rideReq := &ridepb.CreateRideRequest{
		Source:      req.Ride.Source,
		Destination: req.Ride.Destination,
		Distance:    req.Ride.Distance,
		Cost:        req.Ride.Cost,
	}
	rideRes, err := s.rideClient.CreateRide(ctx, rideReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create ride: %v", err)
	}

	booking, err := s.repo.Create(ctx, req.UserId, rideRes.RideId)
	if err != nil {
		log.Printf("failed to insert booking: %v", err)
		return nil, err
	}

	return &pb.Booking{
		BookingId: booking.ID,
		UserId:    booking.UserID,
		RideId:    booking.RideID,
		Time:      booking.Time,
	}, nil
}

func (s *BookingServer) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.BookingDetails, error) {
	booking, err := s.repo.GetByID(ctx, req.BookingId)
	if err != nil {
		return nil, fmt.Errorf("booking not found: %v", err)
	}

	userRes, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{UserId: booking.UserID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	rideRes, err := s.rideClient.GetRide(ctx, &ridepb.GetRideRequest{RideId: booking.RideID})
	if err != nil {
		return nil, fmt.Errorf("failed to get ride: %v", err)
	}

	return &pb.BookingDetails{
		Name:        userRes.Name,
		Source:      rideRes.Source,
		Destination: rideRes.Destination,
		Distance:    rideRes.Distance,
		Cost:        rideRes.Cost,
		Time:        booking.Time,
	}, nil
}
