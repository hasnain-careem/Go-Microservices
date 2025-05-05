package server

import (
	"context"
	"log"

	pb "ride-service/pb/proto/ride"
	"ride-service/repository"
)

type RideServer struct {
	pb.UnimplementedRideServiceServer
	repo repository.RideRepository
}

func NewRideServer(repo repository.RideRepository) *RideServer {
	return &RideServer{repo: repo}
}

func (s *RideServer) CreateRide(ctx context.Context, req *pb.CreateRideRequest) (*pb.CreateRideResponse, error) {
	rideID, err := s.repo.Create(ctx, req.Source, req.Destination, req.Distance, req.Cost)
	if err != nil {
		log.Printf("failed to create ride: %v", err)
		return nil, err
	}

	return &pb.CreateRideResponse{
		RideId: rideID,
	}, nil
}

func (s *RideServer) GetRide(ctx context.Context, req *pb.GetRideRequest) (*pb.Ride, error) {
	ride, err := s.repo.GetByID(ctx, req.RideId)
	if err != nil {
		log.Printf("failed to get ride: %v", err)
		return nil, err
	}

	return &pb.Ride{
		RideId:      ride.ID,
		Source:      ride.Source,
		Destination: ride.Destination,
		Distance:    ride.Distance,
		Cost:        ride.Cost,
	}, nil
}

func (s *RideServer) UpdateRide(ctx context.Context, req *pb.UpdateRideRequest) (*pb.UpdateRideResponse, error) {
	r := req.GetRide()
	message, err := s.repo.Update(ctx, req.RideId, r.Source, r.Destination, r.Distance, r.Cost)
	if err != nil {
		log.Printf("failed to update ride: %v", err)
		return nil, err
	}

	return &pb.UpdateRideResponse{
		Message: message,
	}, nil
}
