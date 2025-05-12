package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"ride-service/config"
	pb "ride-service/pb/proto/ride"
	"ride-service/repository"
	"ride-service/server"

	"github.com/hasnain-zafar/go-microservices/common/metrics"
)

func main() {
	// Initialize Prometheus metrics
	metrics.Init()

	// Start metrics HTTP server in a goroutine
	go startMetricsServer("ride-service", 2113)

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

	fmt.Println("✅ Connected to rides_db successfully")

	rideRepo := repository.NewPostgresRideRepository(db)

	rideServer := server.NewRideServer(rideRepo)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("❌ Failed to listen on port 50052: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRideServiceServer(grpcServer, rideServer)

	reflection.Register(grpcServer)

	fmt.Println("🚀 RideService gRPC server listening on :50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}

func startMetricsServer(serviceName string, port int) {
	fmt.Printf("📊 Metrics server for %s starting on :%d\n", serviceName, port)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("❌ Failed to start metrics server: %v", err)
	}
}
