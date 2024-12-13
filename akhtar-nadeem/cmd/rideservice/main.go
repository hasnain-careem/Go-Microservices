package main

import (
	"log"
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
)

func main() {
	logging.Init()
	db, err := database.DBInit()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	metrics.StartMetricsServer(":9092")

	ridesRepo := repositories.NewRidesRepository(db)
	ridesService := services.NewRidesService(ridesRepo, logging.Logger)

	port := os.Getenv("RIDE_SERVICE_PORT")
	if port == "" {
		port = "50053"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterRidesServiceServer(grpcServer, ridesService)
	logging.Logger.Infof("RidesService listening on port %s", port)

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
	logging.Logger.Info("RidesService: Shutting down server...")

	// Gracefully stop the gRPC server
	grpcServer.GracefulStop()
	logging.Logger.Info("RidesService: Server stopped")
}
