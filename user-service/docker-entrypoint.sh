#!/bin/sh
set -e

echo "Starting User Service..."
echo "Connecting to database at $DB_HOST:$DB_PORT"

# Execute the binary
./user-service
