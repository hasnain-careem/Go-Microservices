package server

import (
	"context"
	"fmt"

	pb "ride-service/pb/proto/ride"
	"ride-service/repository"

	"github.com/hasnain-zafar/go-microservices/common/errors"
	"github.com/hasnain-zafar/go-microservices/common/logger"
)

type RideServer struct {
	pb.UnimplementedRideServiceServer
	repo         repository.RideRepository
	logger       *logger.Logger
	errorHandler *errors.ErrorHandler
}

func NewRideServer(repo repository.RideRepository) *RideServer {
	log := logger.NewLogger("ride-service")
	return &RideServer{
		repo:         repo,
		logger:       log,
		errorHandler: errors.NewErrorHandler(log),
	}
}

func (s *RideServer) CreateRide(ctx context.Context, req *pb.CreateRideRequest) (*pb.CreateRideResponse, error) {
	s.logger.LogRequest("CreateRide", req)

	if err := validateCreateRideRequest(req); err != nil {
		return nil, s.errorHandler.HandleInvalidArgument("invalid ride request", err)
	}

	rideID, err := s.repo.Create(ctx, req.Source, req.Destination, req.Distance, req.Cost)
	if err != nil {
		return nil, s.errorHandler.HandleDatabaseError("failed to create ride", err)
	}

	res := &pb.CreateRideResponse{
		RideId: rideID,
	}

	s.logger.LogResponse("CreateRide", res)

	return res, nil
}

func (s *RideServer) GetRide(ctx context.Context, req *pb.GetRideRequest) (*pb.Ride, error) {
	s.logger.LogRequest("GetRide", req)

	if req.GetRideId() <= 0 {
		return nil, s.errorHandler.HandleInvalidArgument("invalid ride ID", fmt.Errorf("ride ID must be positive"))
	}

	ride, err := s.repo.GetByID(ctx, req.RideId)
	if err != nil {
		if err.Error() == "ride not found" {
			return nil, s.errorHandler.HandleNotFound("ride not found", err)
		}
		return nil, s.errorHandler.HandleDatabaseError("failed to get ride", err)
	}

	res := &pb.Ride{
		RideId:      ride.ID,
		Source:      ride.Source,
		Destination: ride.Destination,
		Distance:    ride.Distance,
		Cost:        ride.Cost,
	}

	s.logger.LogResponse("GetRide", res)

	return res, nil
}

func (s *RideServer) UpdateRide(ctx context.Context, req *pb.UpdateRideRequest) (*pb.UpdateRideResponse, error) {
	s.logger.LogRequest("UpdateRide", req)

	if req.GetRideId() <= 0 {
		return nil, s.errorHandler.HandleInvalidArgument("invalid ride ID", fmt.Errorf("ride ID must be positive"))
	}

	if req.GetRide() == nil {
		return nil, s.errorHandler.HandleInvalidArgument("missing ride details", fmt.Errorf("ride details are required"))
	}

	r := req.GetRide()
	if err := validateRideDetails(r); err != nil {
		return nil, s.errorHandler.HandleInvalidArgument("invalid ride details", err)
	}

	message, err := s.repo.Update(ctx, req.RideId, r.Source, r.Destination, r.Distance, r.Cost)
	if err != nil {
		return nil, s.errorHandler.HandleDatabaseError("failed to update ride", err)
	}

	res := &pb.UpdateRideResponse{
		Message: message,
	}

	s.logger.LogResponse("UpdateRide", res)

	return res, nil
}

func validateCreateRideRequest(req *pb.CreateRideRequest) error {
	if req.Source == "" {
		return fmt.Errorf("source cannot be empty")
	}
	if req.Destination == "" {
		return fmt.Errorf("destination cannot be empty")
	}
	if req.Distance <= 0 {
		return fmt.Errorf("distance must be positive")
	}
	if req.Cost <= 0 {
		return fmt.Errorf("cost must be positive")
	}
	return nil
}

func validateRideDetails(ride *pb.Ride) error {
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
