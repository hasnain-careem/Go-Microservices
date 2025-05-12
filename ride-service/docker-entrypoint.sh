#!/bin/sh
set -e

echo "Starting Ride Service..."
echo "Connecting to database at $DB_HOST:$DB_PORT"

# Execute the binary
./ride-service
