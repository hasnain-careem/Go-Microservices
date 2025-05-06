#!/bin/sh
set -e

echo "Starting Booking Service..."
echo "Connecting to database at $DB_HOST:$DB_PORT"
echo "Depending on User Service and Ride Service"

# Execute the binary
./booking-service
