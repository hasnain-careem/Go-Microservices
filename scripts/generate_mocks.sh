#!/bin/bash

# Exit on error
set -e

# Ensure mockery is installed
if ! command -v mockery &> /dev/null; then
    echo "Mockery is not installed. Installing..."
    go install github.com/vektra/mockery/v3@latest
fi

cd $(dirname $0)/..
PROJECT_ROOT=$(pwd)

echo "Generating mocks for user-service repositories..."
cd $PROJECT_ROOT/user-service
mockery --name=UserRepository --dir=repository --output=repository/mocks --outpkg=mocks

echo "Generating mocks for user-service client..."
cd $PROJECT_ROOT/user-service
mockery --name=UserServiceClient --dir=pb/proto/user --output=pb/proto/user/mocks --outpkg=mocks

echo "Generating mocks for ride-service repositories..."
cd $PROJECT_ROOT/ride-service
mockery --name=RideRepository --dir=repository --output=repository/mocks --outpkg=mocks

echo "Generating mocks for ride-service client..."
cd $PROJECT_ROOT/ride-service
mockery --name=RideServiceClient --dir=pb/proto/ride --output=pb/proto/ride/mocks --outpkg=mocks

echo "Generating mocks for booking-service repositories..."
cd $PROJECT_ROOT/booking-service
mockery --name=BookingRepository --dir=repository --output=repository/mocks --outpkg=mocks

echo "All mocks generated successfully!"
