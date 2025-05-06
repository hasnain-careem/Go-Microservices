package server

import (
	pb "booking-service/pb/proto/booking"
	"booking-service/repository"
	"context"
	"fmt"

	ridepb "ride-service/pb/proto/ride"
	userpb "user-service/pb/proto/user"

	"github.com/hasnain-zafar/go-microservices/common/errors"
	"github.com/hasnain-zafar/go-microservices/common/logger"
)

type BookingServer struct {
	pb.UnimplementedBookingServiceServer
	repo         repository.BookingRepository
	userClient   userpb.UserServiceClient
	rideClient   ridepb.RideServiceClient
	logger       *logger.Logger
	errorHandler *errors.ErrorHandler
}

func NewBookingServer(
	repo repository.BookingRepository,
	userClient userpb.UserServiceClient,
	rideClient ridepb.RideServiceClient,
) *BookingServer {
	log := logger.NewLogger("booking-service")
	return &BookingServer{
		repo:         repo,
		userClient:   userClient,
		rideClient:   rideClient,
		logger:       log,
		errorHandler: errors.NewErrorHandler(log),
	}
}

func (s *BookingServer) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.Booking, error) {
	s.logger.LogRequest("CreateBooking", req)

	if err := validateCreateBookingRequest(req); err != nil {
		return nil, s.errorHandler.HandleInvalidArgument("invalid booking request", err)
	}

	_, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{UserId: req.UserId})
	if err != nil {
		s.logger.Error("failed to get user", "error", err, "user_id", req.UserId)
		logger.IncrementNetworkErrorCount()
		return nil, s.errorHandler.HandleNetworkError("failed to verify user", err)
	}

	rideReq := &ridepb.CreateRideRequest{
		Source:      req.Ride.Source,
		Destination: req.Ride.Destination,
		Distance:    req.Ride.Distance,
		Cost:        req.Ride.Cost,
	}

	rideRes, err := s.rideClient.CreateRide(ctx, rideReq)
	if err != nil {
		s.logger.Error("failed to create ride", "error", err)
		logger.IncrementNetworkErrorCount()
		return nil, s.errorHandler.HandleNetworkError("failed to create ride", err)
	}

	booking, err := s.repo.Create(ctx, req.UserId, rideRes.RideId)
	if err != nil {
		return nil, s.errorHandler.HandleDatabaseError("failed to create booking", err)
	}

	res := &pb.Booking{
		BookingId: booking.ID,
		UserId:    booking.UserID,
		RideId:    booking.RideID,
		Time:      booking.Time,
	}

	s.logger.LogResponse("CreateBooking", res)

	return res, nil
}

func (s *BookingServer) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.BookingDetails, error) {
	s.logger.LogRequest("GetBooking", req)

	if req.GetBookingId() <= 0 {
		return nil, s.errorHandler.HandleInvalidArgument("invalid booking ID", fmt.Errorf("booking ID must be positive"))
	}

	booking, err := s.repo.GetByID(ctx, req.BookingId)
	if err != nil {
		if err.Error() == "booking not found" {
			return nil, s.errorHandler.HandleNotFound("booking not found", err)
		}
		return nil, s.errorHandler.HandleDatabaseError("failed to get booking", err)
	}

	userRes, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{UserId: booking.UserID})
	if err != nil {
		s.logger.Error("failed to get user details", "error", err, "user_id", booking.UserID)
		logger.IncrementNetworkErrorCount()
		return nil, s.errorHandler.HandleNetworkError("failed to get user details", err)
	}

	rideRes, err := s.rideClient.GetRide(ctx, &ridepb.GetRideRequest{RideId: booking.RideID})
	if err != nil {
		s.logger.Error("failed to get ride details", "error", err, "ride_id", booking.RideID)
		logger.IncrementNetworkErrorCount()
		return nil, s.errorHandler.HandleNetworkError("failed to get ride details", err)
	}

	res := &pb.BookingDetails{
		Name:        userRes.Name,
		Source:      rideRes.Source,
		Destination: rideRes.Destination,
		Distance:    rideRes.Distance,
		Cost:        rideRes.Cost,
		Time:        booking.Time,
	}

	s.logger.LogResponse("GetBooking", res)

	return res, nil
}

func validateCreateBookingRequest(req *pb.CreateBookingRequest) error {
	if req.UserId <= 0 {
		return fmt.Errorf("user ID must be positive")
	}

	if req.Ride == nil {
		return fmt.Errorf("ride details are required")
	}

	ride := req.Ride
	if ride.Source == "" {
		return fmt.Errorf("source cannot be empty")
	}
	if ride.Destination == "" {
		return fmt.Errorf("destination cannot be empty")
	}
	if ride.Distance <= 0 {
		return fmt.Errorf("distance must be positive")
	}
	if ride.Cost <= 0 {
		return fmt.Errorf("cost must be positive")
	}

	return nil
}
