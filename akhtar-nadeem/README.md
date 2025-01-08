# Careem Assignment

This repository contains a microservices-based application built with Go, using gRPC for communication and PostgreSQL for data storage. The system consists of three services:

- **User Service (`userservice`)**: Manages user details
- **Booking Service (`bookingservice`)**: Handles bookings and references user and ride data
- **Ride Service (`rideservice`)**: Manages ride information

## Prerequisites

- **Go (1.23 or higher)**
- **Docker & Docker Compose**
- **protoc** (Protocol Buffers compiler)
- **protoc-gen-go** and **protoc-gen-go-grpc** plugins

## Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/akhtarCareem/golang-assignment.git
   cd golang-assignment
   ```



2. **Build and start services**: 
   ```bash
   docker-compose build
   docker-compose up -d
   ```

3. **Check logs**:
   ```bash
   docker-compose logs -f
   ```

4. **Service Ports**:

- User Service: gRPC on 50051, metrics on 9090
- Booking Service: gRPC on 50052, metrics on 9091
- Ride Service: gRPC on 50053, metrics on 9092
- PostgreSQL: 5432

5. **Run tests**:
   ```bash
   go test ./...
   ```

Troubleshooting
•	If builds fail, ensure proto/ directory is included and run protoc again.
•	Check environment variables and .env files for database credentials.
•	Run docker-compose logs for more details.

Happy coding!

