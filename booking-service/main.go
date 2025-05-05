package main

import (
	"booking-service/config"
	"booking-service/pb/proto/booking"
	"booking-service/server"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"


	userpb "user-service/pb/proto/user"
	ridepb "ride-service/pb/proto/ride"
)

func main() {
	// Load env-based config
	cfg := config.Load()

	// Connect to local bookings_db
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

	// Connect to user-service
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Failed to connect to user-service: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	// Connect to ride-service
	rideConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Failed to connect to ride-service: %v", err)
	}
	defer rideConn.Close()
	rideClient := ridepb.NewRideServiceClient(rideConn)

	// Start gRPC server for BookingService
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("❌ Failed to listen on port 50053: %v", err)
	}

	grpcServer := grpc.NewServer()
	bookingServer := &server.BookingServer{
		DB:         db,
		UserClient: userClient,
		RideClient: rideClient,
	}
	pb.RegisterBookingServiceServer(grpcServer, bookingServer)

	// Enable reflection for gRPC server
	reflection.Register(grpcServer)
	
	fmt.Println("🚀 BookingService gRPC server listening on :50053")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}
