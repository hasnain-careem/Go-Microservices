# Go Microservices Architecture for Ride Booking System

This project demonstrates a microservices architecture for a ride booking system implemented in Go using gRPC for service communication and PostgreSQL for data persistence. The architecture includes logging, dependency injection, Prometheus metrics, and comprehensive unit tests.

## System Architecture

The system consists of three core microservices:

1. **UserService (Port 50051)**: Manages user information
2. **RideService (Port 50052)**: Handles ride details
3. **BookingService (Port 50053)**: Coordinates bookings between users and rides

### Service Capabilities

| Service | Capabilities |
|---------|--------------|
| **UserService** | Create, retrieve, and delete users |
| **RideService** | Create, retrieve, and update rides |
| **BookingService** | Create bookings and fetch detailed booking information (using data from UserService and RideService) |

## Technical Implementation

### Key Features

- **gRPC Communication**: Services communicate using gRPC for efficient, type-safe inter-service calls
- **Dependency Injection**: Clean architecture with dependency injection for flexible component management
- **Structured Logging**: Comprehensive logging system to track service operations
- **Prometheus Metrics**: Performance monitoring and metrics collection
- **Unit Testing**: Thorough test coverage with mocking of dependencies
- **PostgreSQL Databases**: Each service has its own dedicated database

## Project Structure

```
go-microservices/
├── common/                   # Shared libraries and utilities
│   ├── errors/               # Common error definitions
│   ├── logger/               # Logging implementation
│   └── metrics/              # Prometheus metrics
│
├── proto/                    # Proto definitions
│   ├── booking/              # Booking service protobuf definitions
│   ├── ride/                 # Ride service protobuf definitions
│   └── user/                 # User service protobuf definitions
│
├── user-service/
│   ├── config/               # Configuration management
│   ├── db/migrations/        # Database schema
│   ├── pb/proto/user/        # Generated protobuf code
│   │   └── mocks/            # Mocks for testing
│   ├── repository/           # Database interactions
│   │   └── mocks/            # Repository mocks
│   ├── server/               # Service implementation
│   ├── .env                  # Environment variables
│   └── main.go               # Entry point
│
├── ride-service/             # Similar structure to user-service
│
├── booking-service/          # Similar structure to user-service
│
└── scripts/                  # Utility scripts
    └── generate_mocks.sh     # Script to generate mocks
```

## Getting Started

### Prerequisites

- Go (1.18+)
- PostgreSQL
- Protocol Buffers compiler (protoc)
- gRPC tools

### Setup PostgreSQL

1. Install PostgreSQL:

```bash
# For macOS
brew install postgresql
brew services start postgresql

# For Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

2. Create the necessary databases:

```bash
createdb users_db
createdb rides_db
createdb bookings_db
```

3. Run the migrations to set up the database schema:

```bash
# Create user table and sample data
psql -d users_db -f user-service/db/migrations/001_create_users_table.sql

# Create ride table and sample data
psql -d rides_db -f ride-service/db/migrations/001_create_rides_table.sql

# Create booking table and sample data
psql -d bookings_db -f booking-service/db/migrations/001_create_bookings_table.sql
```

### Environment Setup

Each service has its own `.env` file containing environment-specific configurations. Update these files according to your local setup:

```bash
# Example .env file for user-service
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=users_db
DB_HOST=localhost
DB_PORT=5432
```

Repeat this for the `ride-service` and `booking-service` directories, adjusting the values as needed.

### Installing Dependencies

For each service, you need to run:

```bash
# From the project root
cd common
go mod download

# For each service directory
cd ../user-service
go mod download

cd ../ride-service
go mod download

cd ../booking-service
go mod download
```

### Running the Services

Start each service in a separate terminal window:

#### 1. User Service

```bash
cd user-service
go run main.go
```

#### 2. Ride Service

```bash
cd ride-service
go run main.go
```

#### 3. Booking Service

```bash
cd booking-service
go run main.go
```

## Testing the Services

### Using gRPC Clients

You can test the gRPC services using `grpcurl`, a command-line tool to interact with gRPC servers:

```bash
# Install grpcurl
# For macOS
brew install grpcurl

# For Linux
# Download from https://github.com/fullstorydev/grpcurl/releases
```

#### Testing UserService

```bash
# List available methods
grpcurl -plaintext localhost:50051 list user.UserService

# Create a user
grpcurl -plaintext -d '{"name": "John Doe"}' localhost:50051 user.UserService/CreateUser

# Get a user (replace 1 with the user_id from the previous command)
grpcurl -plaintext -d '{"user_id": 1}' localhost:50051 user.UserService/GetUser

# Delete a user
grpcurl -plaintext -d '{"user_id": 1}' localhost:50051 user.UserService/DeleteUser
```

#### Testing RideService

```bash
# List available methods
grpcurl -plaintext localhost:50052 list ride.RideService

# Create a ride
grpcurl -plaintext -d '{"source": "New York", "destination": "Boston", "distance": 215, "cost": 120}' localhost:50052 ride.RideService/CreateRide

# Get a ride (replace 1 with the ride_id from the previous command)
grpcurl -plaintext -d '{"ride_id": 1}' localhost:50052 ride.RideService/GetRide

# Update a ride
grpcurl -plaintext -d '{"ride_id": 1, "ride": {"source": "New York", "destination": "Philadelphia", "distance": 150, "cost": 80}}' localhost:50052 ride.RideService/UpdateRide
```

#### Testing BookingService

```bash
# List available methods
grpcurl -plaintext localhost:50053 list booking.BookingService

# Create a booking (use a valid user_id)
grpcurl -plaintext -d '{"user_id": 2, "ride": {"source": "Chicago", "destination": "Detroit", "distance": 280, "cost": 150}}' localhost:50053 booking.BookingService/CreateBooking

# Get booking details (replace 1 with the booking_id from the previous command)
grpcurl -plaintext -d '{"booking_id": 1}' localhost:50053 booking.BookingService/GetBooking
```

### Unit Testing

The project includes comprehensive unit tests for each service. To run the tests:

```bash
# Run tests for user-service
cd user-service
go test ./... -v

# Run tests for ride-service
cd ../ride-service
go test ./... -v

# Run tests for booking-service
cd ../booking-service
go test ./... -v
```

### Test Coverage

To check the test coverage for a service:

```bash
cd user-service
go test ./... -cover

# For more detailed coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Mock Generation

The project uses mocks for testing service interfaces. You can generate mocks using the provided script or manually:

#### Using the Script

```bash
# Make the script executable if needed
chmod +x scripts/generate_mocks.sh

# Run the script
./scripts/generate_mocks.sh
```

#### Manually Generating Mocks

1. Install mockgen:

```bash
go install github.com/golang/mock/mockgen@latest
```

2. Generate mocks for a service (example for UserService):

```bash
# For repository
cd user-service
mockgen -source=repository/user_repository.go -destination=repository/mocks/UserRepository.go -package=mocks

# For gRPC service client
mockgen -source=pb/proto/user/user_grpc.pb.go -destination=pb/proto/user/mocks/UserServiceClient.go -package=mocks
```

## Monitoring with Prometheus

Each service exposes metrics on a dedicated port that can be scraped by Prometheus:

- UserService: http://localhost:9091/metrics
- RideService: http://localhost:9092/metrics  
- BookingService: http://localhost:9093/metrics

To view these metrics, you can:

1. Install Prometheus:

```bash
# For macOS
brew install prometheus

# For Ubuntu/Debian
wget https://github.com/prometheus/prometheus/releases/download/v2.38.0/prometheus-2.38.0.linux-amd64.tar.gz
tar xvfz prometheus-2.38.0.linux-amd64.tar.gz
```

2. Configure Prometheus to scrape your services by creating a `prometheus.yml` file:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'user-service'
    static_configs:
      - targets: ['localhost:9091']
  
  - job_name: 'ride-service'
    static_configs:
      - targets: ['localhost:9092']
  
  - job_name: 'booking-service'
    static_configs:
      - targets: ['localhost:9093']
```

3. Start Prometheus:

```bash
prometheus --config.file=prometheus.yml
```

4. Access the Prometheus UI at http://localhost:9090

## Load Testing

You can use `hey` for HTTP-based load testing:

```bash
# Install hey
go install github.com/rakyll/hey@latest

# Example load test for booking service
hey -n 100 -c 10 -m POST -d '{"user_id": 2, "ride": {"source": "Chicago", "destination": "Detroit", "distance": 280, "cost": 150}}' http://localhost:50053/bookings
```

## Dependency Injection

The project uses a clean architecture with dependency injection to ensure components are loosely coupled:

1. **Repository Layer**: Handles database operations
2. **Service Layer**: Contains business logic  
3. **Handler Layer**: Processes gRPC requests and responses

This approach makes testing easier and components more reusable.

## API Definitions

### UserService

- **CreateUser**: Creates a new user
  - Request: `{"name": "string"}`
  - Response: `{"user_id": int}`

- **GetUser**: Retrieves a user by ID
  - Request: `{"user_id": int}`
  - Response: `{"name": "string"}`

- **DeleteUser**: Deletes a user by ID
  - Request: `{"user_id": int}`
  - Response: `{"message": "string"}`

### RideService

- **CreateRide**: Creates a new ride
  - Request: `{"source": "string", "destination": "string", "distance": int, "cost": int}`
  - Response: `{"ride_id": int}`

- **GetRide**: Retrieves ride details
  - Request: `{"ride_id": int}`
  - Response: `{"ride_id": int, "source": "string", "destination": "string", "distance": int, "cost": int}`

- **UpdateRide**: Updates an existing ride
  - Request: `{"ride_id": int, "ride": {"source": "string", "destination": "string", "distance": int, "cost": int}}`
  - Response: `{"message": "string"}`

### BookingService

- **CreateBooking**: Creates a new booking
  - Request: `{"user_id": int, "ride": {"source": "string", "destination": "string", "distance": int, "cost": int}}`
  - Response: `{"booking_id": int, "user_id": int, "ride_id": int, "time": "string"}`

- **GetBooking**: Retrieves booking details with user and ride information
  - Request: `{"booking_id": int}`
  - Response: `{"name": "string", "source": "string", "destination": "string", "distance": int, "cost": int, "time": "string"}`

## Troubleshooting

### Common Issues

1. **Database Connection Problems**:
   - Verify PostgreSQL is running: `pg_isready`
   - Check credentials in `.env` files
   - Ensure database exists: `psql -l`

2. **gRPC Service Not Starting**:
   - Check if port is already in use: `lsof -i :<port>`
   - Verify environment variables are loaded properly

3. **Dependency Issues**:
   - Run `go mod tidy` to clean up and update dependencies
   - Check that all required modules are downloaded

4. **Proto Generation Issues**:
   - Ensure you have the latest version of protoc installed
   - Check that all proto imports are correctly specified

### Logs

Each service writes logs that can help diagnose issues. Check the console output or configured log files.
