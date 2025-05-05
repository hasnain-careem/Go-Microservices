package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	pb "ride-service/pb/proto/ride"
)

type RideServer struct {
	pb.UnimplementedRideServiceServer
	DB *sql.DB
}

func (s *RideServer) UpdateRide(ctx context.Context, req *pb.UpdateRideRequest) (*pb.UpdateRideResponse, error) {
	query := `UPDATE rides SET source = $1, destination = $2, distance = $3, cost = $4 WHERE ride_id = $5`

	r := req.GetRide()
	_, err := s.DB.ExecContext(ctx, query, r.Source, r.Destination, r.Distance, r.Cost, req.RideId)
	if err != nil {
		log.Printf("failed to update ride: %v", err)
		return nil, err
	}

	return &pb.UpdateRideResponse{
		Message: fmt.Sprintf("Ride %d updated successfully", req.RideId),
	}, nil
}

// Optional: add GetRide to support BookingService.GetBooking
func (s *RideServer) GetRide(ctx context.Context, req *pb.GetRideRequest) (*pb.Ride, error) {
	query := `SELECT source, destination, distance, cost FROM rides WHERE ride_id = $1`
	var source, destination string
	var distance, cost int32

	err := s.DB.QueryRowContext(ctx, query, req.RideId).Scan(&source, &destination, &distance, &cost)
	if err != nil {
		return nil, fmt.Errorf("ride not found: %v", err)
	}

	return &pb.Ride{
		RideId:      req.RideId,
		Source:      source,
		Destination: destination,
		Distance:    distance,
		Cost:        cost,
	}, nil
}

func (s *RideServer) CreateRide(ctx context.Context, req *pb.CreateRideRequest) (*pb.CreateRideResponse, error) {
	query := `INSERT INTO rides (source, destination, distance, cost) VALUES ($1, $2, $3, $4) RETURNING ride_id`

	var rideID int32
	err := s.DB.QueryRowContext(ctx, query, req.Source, req.Destination, req.Distance, req.Cost).Scan(&rideID)
	if err != nil {
		log.Printf("failed to create ride: %v", err)
		return nil, err
	}

	return &pb.CreateRideResponse{
		RideId: rideID,
	}, nil
}
