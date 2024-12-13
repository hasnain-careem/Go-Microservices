package services

import (
	"context"

	"github.com/akhtarCareem/golang-assignment/internal/repositories"
	"github.com/akhtarCareem/golang-assignment/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RidesServiceServer struct {
	proto.UnimplementedRidesServiceServer
	repo repositories.RidesRepository
	log  *logrus.Logger
}

func NewRidesService(repo repositories.RidesRepository, log *logrus.Logger) *RidesServiceServer {
	return &RidesServiceServer{repo: repo, log: log}
}

func (r *RidesServiceServer) UpdateRide(ctx context.Context, req *proto.UpdateRideRequest) (*proto.UpdateRideResponse, error) {
	r.log.WithField("ride_id", req.RideId).Info("UpdateRide called")
	err := r.repo.UpdateRide(req.RideId, req.Ride.Source, req.Ride.Destination, req.Ride.Distance, req.Ride.Cost)
	if err != nil {
		r.log.WithError(err).Error("failed to update ride")
		return nil, status.Error(codes.NotFound, "ride not found")
	}
	return &proto.UpdateRideResponse{Message: "ride updated successfully"}, nil
}
