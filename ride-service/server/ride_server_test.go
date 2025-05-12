package server

import (
	"context"
	"errors"
	"testing"

	pb "ride-service/pb/proto/ride"
	"ride-service/repository"
	"ride-service/repository/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateRide_Success(t *testing.T) {
	// Setup
	mockRepo := new(mocks.RideRepository)
	rideServer := NewRideServer(mockRepo)

	ctx := context.Background()
	req := &pb.CreateRideRequest{
		Source:      "New York",
		Destination: "Boston",
		Distance:    200,
		Cost:        150,
	}

	// Expectations
	mockRepo.On("Create", ctx, "New York", "Boston", int32(200), int32(150)).Return(int32(1), nil)

	// Action
	resp, err := rideServer.CreateRide(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.RideId)
	mockRepo.AssertExpectations(t)
}

func TestCreateRide_InvalidRequest(t *testing.T) {
	// Create a set of test cases for different validation failures
	testCases := []struct {
		name string
		req  *pb.CreateRideRequest
	}{
		{
			name: "Empty Source",
			req: &pb.CreateRideRequest{
				Source:      "",
				Destination: "Boston",
				Distance:    200,
				Cost:        150,
			},
		},
		{
			name: "Empty Destination",
			req: &pb.CreateRideRequest{
				Source:      "New York",
				Destination: "",
				Distance:    200,
				Cost:        150,
			},
		},
		{
			name: "Zero Distance",
			req: &pb.CreateRideRequest{
				Source:      "New York",
				Destination: "Boston",
				Distance:    0,
				Cost:        150,
			},
		},
		{
			name: "Negative Distance",
			req: &pb.CreateRideRequest{
				Source:      "New York",
				Destination: "Boston",
				Distance:    -10,
				Cost:        150,
			},
		},
		{
			name: "Zero Cost",
			req: &pb.CreateRideRequest{
				Source:      "New York",
				Destination: "Boston",
				Distance:    200,
				Cost:        0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mocks.RideRepository)
			rideServer := NewRideServer(mockRepo)

			// Action
			resp, err := rideServer.CreateRide(context.Background(), tc.req)

			// Assertions
			assert.Error(t, err)
			assert.Nil(t, resp)
			// Repository should not be called when validation fails
			mockRepo.AssertNotCalled(t, "Create")
		})
	}
}

func TestCreateRide_RepositoryError(t *testing.T) {
	// Setup
	mockRepo := new(mocks.RideRepository)
	rideServer := NewRideServer(mockRepo)

	ctx := context.Background()
	req := &pb.CreateRideRequest{
		Source:      "New York",
		Destination: "Boston",
		Distance:    200,
		Cost:        150,
	}

	// Expectations
	mockRepo.On("Create", ctx, "New York", "Boston", int32(200), int32(150)).Return(int32(0), errors.New("database error"))

	// Action
	resp, err := rideServer.CreateRide(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestGetRide_Success(t *testing.T) {
	// Setup
	mockRepo := new(mocks.RideRepository)
	rideServer := NewRideServer(mockRepo)

	ctx := context.Background()
	req := &pb.GetRideRequest{RideId: 1}

	// Mock repository response
	mockRide := &repository.Ride{
		ID:          1,
		Source:      "New York",
		Destination: "Boston",
		Distance:    200,
		Cost:        150,
	}

	// Expectations
	mockRepo.On("GetByID", ctx, int32(1)).Return(mockRide, nil)

	// Action
	resp, err := rideServer.GetRide(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.RideId)
	assert.Equal(t, "New York", resp.Source)
	assert.Equal(t, "Boston", resp.Destination)
	assert.Equal(t, int32(200), resp.Distance)
	assert.Equal(t, int32(150), resp.Cost)
	mockRepo.AssertExpectations(t)
}

func TestGetRide_InvalidId(t *testing.T) {
	// Setup
	mockRepo := new(mocks.RideRepository)
	rideServer := NewRideServer(mockRepo)

	ctx := context.Background()
	req := &pb.GetRideRequest{RideId: 0}

	// Action
	resp, err := rideServer.GetRide(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	// Repository should not be called when validation fails
	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestGetRide_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(mocks.RideRepository)
	rideServer := NewRideServer(mockRepo)

	ctx := context.Background()
	req := &pb.GetRideRequest{RideId: 999}

	// Expectations
	mockRepo.On("GetByID", ctx, int32(999)).Return(nil, errors.New("ride not found"))

	// Action
	resp, err := rideServer.GetRide(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRide_Success(t *testing.T) {
	// Setup
	mockRepo := new(mocks.RideRepository)
	rideServer := NewRideServer(mockRepo)

	ctx := context.Background()
	req := &pb.UpdateRideRequest{
		RideId: 1,
		Ride: &pb.Ride{
			Source:      "New York",
			Destination: "Boston",
			Distance:    200,
			Cost:        150,
		},
	}
	expectedMsg := "Ride 1 updated successfully"

	// Expectations
	mockRepo.On("Update", ctx, int32(1), "New York", "Boston", int32(200), int32(150)).Return(expectedMsg, nil)

	// Action
	resp, err := rideServer.UpdateRide(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedMsg, resp.Message)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRide_InvalidRequest(t *testing.T) {
	// Create a set of test cases for different validation failures
	testCases := []struct {
		name string
		req  *pb.UpdateRideRequest
	}{
		{
			name: "Invalid Ride ID",
			req: &pb.UpdateRideRequest{
				RideId: 0,
				Ride: &pb.Ride{
					Source:      "New York",
					Destination: "Boston",
					Distance:    200,
					Cost:        150,
				},
			},
		},
		{
			name: "Missing Ride Details",
			req: &pb.UpdateRideRequest{
				RideId: 1,
				Ride:   nil,
			},
		},
		{
			name: "Empty Source",
			req: &pb.UpdateRideRequest{
				RideId: 1,
				Ride: &pb.Ride{
					Source:      "",
					Destination: "Boston",
					Distance:    200,
					Cost:        150,
				},
			},
		},
		{
			name: "Empty Destination",
			req: &pb.UpdateRideRequest{
				RideId: 1,
				Ride: &pb.Ride{
					Source:      "New York",
					Destination: "",
					Distance:    200,
					Cost:        150,
				},
			},
		},
		{
			name: "Zero Distance",
			req: &pb.UpdateRideRequest{
				RideId: 1,
				Ride: &pb.Ride{
					Source:      "New York",
					Destination: "Boston",
					Distance:    0,
					Cost:        150,
				},
			},
		},
		{
			name: "Zero Cost",
			req: &pb.UpdateRideRequest{
				RideId: 1,
				Ride: &pb.Ride{
					Source:      "New York",
					Destination: "Boston",
					Distance:    200,
					Cost:        0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(mocks.RideRepository)
			rideServer := NewRideServer(mockRepo)

			// Action
			resp, err := rideServer.UpdateRide(context.Background(), tc.req)

			// Assertions
			assert.Error(t, err)
			assert.Nil(t, resp)
			// Repository should not be called when validation fails
			mockRepo.AssertNotCalled(t, "Update")
		})
	}
}

func TestUpdateRide_RepositoryError(t *testing.T) {
	// Setup
	mockRepo := new(mocks.RideRepository)
	rideServer := NewRideServer(mockRepo)

	ctx := context.Background()
	req := &pb.UpdateRideRequest{
		RideId: 1,
		Ride: &pb.Ride{
			Source:      "New York",
			Destination: "Boston",
			Distance:    200,
			Cost:        150,
		},
	}

	// Expectations
	mockRepo.On("Update", ctx, int32(1), "New York", "Boston", int32(200), int32(150)).Return("", errors.New("database error"))

	// Action
	resp, err := rideServer.UpdateRide(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}
