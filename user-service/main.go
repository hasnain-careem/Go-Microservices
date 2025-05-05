package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"user-service/config"
	pb "user-service/pb/proto/user"
	"user-service/server"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("❌ Could not connect to DB: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ DB not reachable: %v", err)
	}

	fmt.Println("✅ Connected to users_db successfully")

	// Initialize gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	userServer := &server.UserServer{
		DB: db,
	}
	pb.RegisterUserServiceServer(grpcServer, userServer)

	// Enable reflection
	reflection.Register(grpcServer)

	fmt.Println("🚀 UserService gRPC server listening on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}
