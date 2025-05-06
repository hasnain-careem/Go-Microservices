# Go Microservices Project

This project demonstrates a microservices architecture built with Go, gRPC, and PostgreSQL. It consists of three services:

- **User Service** - Manages user information
- **Ride Service** - Handles ride details and pricing
- **Booking Service** - Coordinates bookings between users and rides

## Architecture

The microservices communicate with each other via gRPC, persist data in PostgreSQL databases, and expose metrics for monitoring with Prometheus.

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │     │                 │
│   User Service  │◄────┤ Booking Service │────►│   Ride Service  │
│  (Port: 50051)  │     │  (Port: 50053)  │     │  (Port: 50052)  │
│                 │     │                 │     │                 │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │     │                 │
│    users_db     │     │   bookings_db   │     │    rides_db     │
│                 │     │                 │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- gRPCurl (for testing gRPC endpoints)

### Running the Services with Docker

1. Clone the repository:

```bash
git clone https://github.com/hasnain-careem/Go-Microservices.git
cd go-microservices
```

2. Start all services using Docker Compose:

```bash
docker-compose up -d
```

This will start:
- Three microservices: user-service, ride-service, and booking-service
- Three PostgreSQL databases: users_db, rides_db, and bookings_db
- Prometheus for metrics collection

3. Verify all services are running:

```bash
docker-compose ps
```

## Monitoring with Prometheus

Prometheus is configured to scrape metrics from all three services:

- User Service metrics: http://localhost:2112/metrics
- Ride Service metrics: http://localhost:2113/metrics
- Booking Service metrics: http://localhost:2114/metrics

Access the Prometheus dashboard at: http://localhost:9090

### Available Metrics

- `grpc_requests_total` - Counter for gRPC requests by service and method
- `app_errors_total` - Counter for errors by service and type

Example Prometheus queries:
- Request rate: `rate(grpc_requests_total[1m])`
- Error rate: `rate(app_errors_total[1m])`
- Success rate: `sum(rate(grpc_requests_total[1m])) - sum(rate(app_errors_total[1m]))`

## Testing with gRPCurl

[gRPCurl](https://github.com/fullstorydev/grpcurl) is a command-line tool that lets you interact with gRPC servers.

### User Service (Port 50051)

List available methods:
```bash
grpcurl -plaintext localhost:50051 list
```

Create a user:
```bash
grpcurl -plaintext -d '{"name": "John Smith"}' localhost:50051 user.UserService/CreateUser
```

Get a user:
```bash
grpcurl -plaintext -d '{"user_id": 1}' localhost:50051 user.UserService/GetUser
```

Delete a user:
```bash
grpcurl -plaintext -d '{"user_id": 1}' localhost:50051 user.UserService/DeleteUser
```

### Ride Service (Port 50052)

List available methods:
```bash
grpcurl -plaintext localhost:50052 list
```

Create a ride:
```bash
grpcurl -plaintext -d '{"source": "New York", "destination": "Boston", "distance": 200, "cost": 150}' localhost:50052 ride.RideService/CreateRide
```

Get a ride:
```bash
grpcurl -plaintext -d '{"ride_id": 1}' localhost:50052 ride.RideService/GetRide
```

Update a ride:
```bash
grpcurl -plaintext -d '{"ride_id": 1, "ride": {"source": "New York", "destination": "Washington DC", "distance": 225, "cost": 175}}' localhost:50052 ride.RideService/UpdateRide
```

### Booking Service (Port 50053)

List available methods:
```bash
grpcurl -plaintext localhost:50053 list
```

Create a booking:
```bash
grpcurl -plaintext -d '{"user_id": 1, "ride": {"source": "Philadelphia", "destination": "Pittsburgh", "distance": 305, "cost": 200}}' localhost:50053 booking.BookingService/CreateBooking
```

Get booking details:
```bash
grpcurl -plaintext -d '{"booking_id": 1}' localhost:50053 booking.BookingService/GetBooking
```

## Project Structure

```
go-microservices/
├── common/              # Shared libraries
│   ├── errors/          # Error handling
│   ├── logger/          # Logging
│   └── metrics/         # Prometheus metrics
├── user-service/        # User microservice
├── ride-service/        # Ride microservice
├── booking-service/     # Booking microservice
├── proto/               # Protocol buffer definitions
├── docker-compose.yml   # Docker Compose configuration
└── scripts/             # Utility scripts
```

## Development

### Generating Mocks for Testing

```bash
./scripts/generate_mocks.sh
```

### Running Tests

```bash
# Run tests for all services
go test ./...

# Run tests for a specific service
cd user-service
go test ./...
```

### Resetting and Rebuilding Docker

To completely reset Docker containers and rebuild the application:

```bash
# Stop and remove all containers
docker-compose down

# Remove all images related to the project (optional)
docker-compose down --rmi all

# Remove all volumes to clear persistent data (optional)
docker-compose down -v

# Rebuild and start services
docker-compose up --build -d
```

For rebuilding specific services only:

```bash
# Rebuild a specific service
docker-compose build user-service

# Restart a specific service
docker-compose up -d --no-deps user-service
```

## Troubleshooting

- If services can't connect to each other, make sure the Docker networks are correctly set up
- Check container logs: `docker-compose logs -f <service-name>`
- Verify database migrations ran correctly by connecting to the databases

## License

This project is licensed under the MIT License - see the LICENSE file for details.
