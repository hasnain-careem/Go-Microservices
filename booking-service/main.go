package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"booking-service/config"
	pb "booking-service/pb/proto/booking"
	"booking-service/repository"
	"booking-service/server"
	ridepb "ride-service/pb/proto/ride"
	userpb "user-service/pb/proto/user"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ Cannot ping DB: %v", err)
	}
	fmt.Println("✅ Connected to bookings_db")

	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Failed to connect to user-service: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	rideConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Failed to connect to ride-service: %v", err)
	}
	defer rideConn.Close()
	rideClient := ridepb.NewRideServiceClient(rideConn)

	bookingRepo := repository.NewPostgresBookingRepository(db)

	bookingServer := server.NewBookingServer(bookingRepo, userClient, rideClient)

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("❌ Failed to listen on port 50053: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, bookingServer)

	reflection.Register(grpcServer)

	fmt.Println("🚀 BookingService gRPC server listening on :50053")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}
