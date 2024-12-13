package services

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/akhtarCareem/golang-assignment/internal/repositories"
	"github.com/akhtarCareem/golang-assignment/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingServiceServer struct {
	proto.UnimplementedBookingServiceServer
	repo repositories.BookingRepository
	log  *logrus.Logger
}

func NewBookingService(repo repositories.BookingRepository, log *logrus.Logger) *BookingServiceServer {
	return &BookingServiceServer{repo: repo, log: log}
}

func (b *BookingServiceServer) CreateBooking(ctx context.Context, req *proto.CreateBookingRequest) (*proto.CreateBookingResponse, error) {
	b.log.WithField("user_id", req.UserId).Info("CreateBooking called")
	bookingID, rideID, bookingTime, err := b.repo.CreateBooking(req.UserId, req.Ride.Source, req.Ride.Destination, req.Ride.Distance, req.Ride.Cost)
	if err != nil {
		b.log.WithError(err).Error("failed to create booking")
		return nil, status.Error(codes.Internal, "failed to create booking")
	}
	return &proto.CreateBookingResponse{
		BookingId: bookingID,
		UserId:    req.UserId,
		RideId:    rideID,
		Time:      toProtoTimestamp(bookingTime),
	}, nil
}

func (b *BookingServiceServer) GetBooking(ctx context.Context, req *proto.GetBookingRequest) (*proto.GetBookingResponse, error) {
	b.log.WithField("booking_id", req.BookingId).Info("GetBooking called")
	name, source, destination, distance, cost, bookingTime, userID, rideID, err := b.repo.GetBookingDetails(req.BookingId)
	if err != nil {
		b.log.WithError(err).Error("failed to get booking")
		return nil, status.Error(codes.NotFound, "booking not found")
	}
	return &proto.GetBookingResponse{
		Name:        name,
		Source:      source,
		Destination: destination,
		Distance:    distance,
		Cost:        cost,
		Time:        toProtoTimestamp(bookingTime),
		UserId:      userID,
		BookingId:   req.BookingId,
		RideId:      rideID,
	}, nil
}

func toProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
