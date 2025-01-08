package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/akhtarCareem/golang-assignment/internal/database"
	"github.com/akhtarCareem/golang-assignment/internal/logging"
	"github.com/akhtarCareem/golang-assignment/internal/metrics"
	"github.com/akhtarCareem/golang-assignment/internal/repositories"
	"github.com/akhtarCareem/golang-assignment/internal/services"
	"github.com/akhtarCareem/golang-assignment/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	logging.Init()
	db, err := database.DBInit()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	metrics.StartMetricsServer(":9090")

	userStore := repositories.NewUserStore(db)
	userService := services.NewUserService(userStore, logging.Logger)

	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, userService)
	logging.Logger.Infof("UserService listening on port %s", port)

	// Graceful Shutdown
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.Logger.Info("UserService: Shutting down server...")

	// Gracefully stop the gRPC server
	grpcServer.GracefulStop()
	logging.Logger.Info("UserService: Server stopped")
}
