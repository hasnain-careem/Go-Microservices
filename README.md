This project demonstrates a microservices architecture for a ride booking system implemented in Go using gRPC for service communication and PostgreSQL for data persistence.

## System Architecture

The system consists of three microservices:

1. **UserService (Port 50051)**: Manages user information
2. **RideService (Port 50052)**: Handles ride details
3. **BookingService (Port 50053)**: Coordinates bookings between users and rides

### Service Capabilities

| Service | Capabilities |
|---------|--------------|
| **UserService** | Create, retrieve, and delete users |
| **RideService** | Create, retrieve, and update rides |
| **BookingService** | Create bookings and fetch detailed booking information (using data from UserService and RideService) |

## Project Structure

```
go-microservices/
├── user-service/
│   ├── config/               # Configuration management
│   ├── pb/proto/user/        # Generated protobuf code
│   ├── server/               # Service implementation
│   ├── db/migrations/        # Database schema
│   ├── .env                  # Environment variables
│   ├── main.go               # Entry point
│   └── go.mod                # Dependencies
│
├── ride-service/
│   ├── config/
│   ├── pb/proto/ride/
│   ├── server/
│   ├── db/migrations/
│   ├── .env
│   ├── main.go
│   └── go.mod
│
├── booking-service/
│   ├── config/
│   ├── pb/proto/booking/
│   ├── server/
│   ├── db/migrations/
│   ├── .env
│   ├── main.go
│   └── go.mod
│
└── proto/                    # Proto definitions
    ├── user/
    ├── ride/
    └── booking/
```

## Getting Started

### Prerequisites

- Go (1.18+)
- PostgreSQL
- Protocol Buffers compiler (protoc)
- gRPC tools
- grpcurl (for testing)

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

Each service has its own .env file. You may need to update the database connection details:

```bash
# Example .env file for user-service
DB_USER=hasnain
DB_PASSWORD=pass123
DB_NAME=users_db
DB_HOST=localhost
DB_PORT=5432
```

Update the credentials in each service's .env file to match your PostgreSQL setup.

### Install Dependencies

For each service, you need to run:

```bash
# From the service directory (user-service, ride-service, booking-service)
go mod download
```

### Running the Services

Start each service in a separate terminal:

#### 1. User Service

```bash
cd user-service
go run main.go
# You should see: "🚀 UserService gRPC server listening on :50051"
```

#### 2. Ride Service

```bash
cd ride-service
go run main.go
# You should see: "🚀 RideService gRPC server listening on :50052"
```

#### 3. Booking Service

```bash
cd booking-service
go run main.go
# You should see: "🚀 BookingService gRPC server listening on :50053"
```

## Testing the Services

You can use `grpcurl` to test the gRPC endpoints. Install it with:

```bash
# For macOS
brew install grpcurl

# For Linux
# Download the binary from https://github.com/fullstorydev/grpcurl/releases
```

### Testing UserService

```bash
# Create a user
grpcurl -plaintext -d '{"name": "John Doe"}' localhost:50051 user.UserService/CreateUser

# Get a user (replace 1 with the user_id from the previous command)
grpcurl -plaintext -d '{"user_id": 1}' localhost:50051 user.UserService/GetUser

# Delete a user
grpcurl -plaintext -d '{"user_id": 1}' localhost:50051 user.UserService/DeleteUser
```

### Testing RideService

```bash
# Create a ride
grpcurl -plaintext -d '{"source": "New York", "destination": "Boston", "distance": 215, "cost": 120}' localhost:50052 ride.RideService/CreateRide

# Get a ride (replace 1 with the ride_id from the previous command)
grpcurl -plaintext -d '{"ride_id": 1}' localhost:50052 ride.RideService/GetRide

# Update a ride
grpcurl -plaintext -d '{"ride_id": 1, "ride": {"source": "New York", "destination": "Philadelphia", "distance": 150, "cost": 80}}' localhost:50052 ride.RideService/UpdateRide
```

### Testing BookingService

```bash
# Create a booking (use a valid user_id)
grpcurl -plaintext -d '{"user_id": 2, "ride": {"source": "Chicago", "destination": "Detroit", "distance": 280, "cost": 150}}' localhost:50053 booking.BookingService/CreateBooking

# Get booking details (replace 1 with the booking_id from the previous command)
grpcurl -plaintext -d '{"booking_id": 1}' localhost:50053 booking.BookingService/GetBooking
```

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

## Architecture Details

- **Inter-service Communication**: Services communicate using gRPC calls
- **Data Storage**: Each service has its own PostgreSQL database
- **Service Discovery**: Hardcoded service addresses (in a production environment, use a service registry)
- **Error Handling**: Each service implements appropriate error handling and logging
