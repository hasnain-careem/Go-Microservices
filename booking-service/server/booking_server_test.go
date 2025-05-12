package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	pb "booking-service/pb/proto/booking"
	"booking-service/repository"
	"booking-service/repository/mocks"
	ridepb "ride-service/pb/proto/ride"
	ridemocks "ride-service/pb/proto/ride/mocks"
	userpb "user-service/pb/proto/user"
	usermocks "user-service/pb/proto/user/mocks"
)

func TestCreateBooking_Success(t *testing.T) {
	// Setup
	mockRepo := new(mocks.BookingRepository)
	mockUserClient := new(usermocks.UserServiceClient)
	mockRideClient := new(ridemocks.RideServiceClient)

	bookingServer := NewBookingServer(mockRepo, mockUserClient, mockRideClient)

	ctx := context.Background()
	req := &pb.CreateBookingRequest{
		UserId: 1,
		Ride: &pb.Ride{
			Source:      "New York",
			Destination: "Boston",
			Distance:    200,
			Cost:        150,
		},
	}

	// Expectations
	mockUserClient.On("GetUser", ctx, &userpb.GetUserRequest{UserId: 1}).
		Return(&userpb.GetUserResponse{Name: "John Doe"}, nil)

	mockRideClient.On("CreateRide", ctx, &ridepb.CreateRideRequest{
		Source:      "New York",
		Destination: "Boston",
		Distance:    200,
		Cost:        150,
	}).Return(&ridepb.CreateRideResponse{RideId: 5}, nil)

	// Mock the booking creation
	mockBooking := &repository.Booking{
		ID:     10,
		UserID: 1,
		RideID: 5,
		Time:   "2023-01-01T12:00:00Z",
	}
	mockRepo.On("Create", ctx, int32(1), int32(5)).Return(mockBooking, nil)

	// Action
	resp, err := bookingServer.CreateBooking(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(10), resp.BookingId)
	assert.Equal(t, int32(1), resp.UserId)
	assert.Equal(t, int32(5), resp.RideId)
	assert.Equal(t, "2023-01-01T12:00:00Z", resp.Time)

	// Verify expectations
	mockUserClient.AssertExpectations(t)
	mockRideClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreateBooking_InvalidRequest(t *testing.T) {
	// Create a set of test cases for different validation failures
	testCases := []struct {
		name string
		req  *pb.CreateBookingRequest
	}{
		{
			name: "Invalid User ID",
			req: &pb.CreateBookingRequest{
				UserId: 0,
				Ride: &pb.Ride{
					Source:      "New York",
					Destination: "Boston",
					Distance:    200,
					Cost:        150,
				},
			},
		},
		{
			name: "Missing Ride",
			req: &pb.CreateBookingRequest{
				UserId: 1,
				Ride:   nil,
			},
		},
		{
			name: "Empty Source",
			req: &pb.CreateBookingRequest{
				UserId: 1,
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
			req: &pb.CreateBookingRequest{
				UserId: 1,
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
			req: &pb.CreateBookingRequest{
				UserId: 1,
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
			req: &pb.CreateBookingRequest{
				UserId: 1,
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
			mockRepo := new(mocks.BookingRepository)
			mockUserClient := new(usermocks.UserServiceClient)
			mockRideClient := new(ridemocks.RideServiceClient)

			bookingServer := NewBookingServer(mockRepo, mockUserClient, mockRideClient)

			// Action
			resp, err := bookingServer.CreateBooking(context.Background(), tc.req)

			// Assertions
			assert.Error(t, err)
			assert.Nil(t, resp)

			// No external services should be called for validation errors
			mockUserClient.AssertNotCalled(t, "GetUser")
			mockRideClient.AssertNotCalled(t, "CreateRide")
			mockRepo.AssertNotCalled(t, "Create")
		})
	}
}

func TestCreateBooking_UserServiceError(t *testing.T) {
	// Setup
	mockRepo := new(mocks.BookingRepository)
	mockUserClient := new(usermocks.UserServiceClient)
	mockRideClient := new(ridemocks.RideServiceClient)

	bookingServer := NewBookingServer(mockRepo, mockUserClient, mockRideClient)

	ctx := context.Background()
	req := &pb.CreateBookingRequest{
		UserId: 999, // Non-existent user
		Ride: &pb.Ride{
			Source:      "New York",
			Destination: "Boston",
			Distance:    200,
			Cost:        150,
		},
	}

	// Expectations
	mockUserClient.On("GetUser", ctx, &userpb.GetUserRequest{UserId: 999}).
		Return(nil, errors.New("user not found"))

	// Action
	resp, err := bookingServer.CreateBooking(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verify expectations
	mockUserClient.AssertExpectations(t)
	mockRideClient.AssertNotCalled(t, "CreateRide")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateBooking_RideServiceError(t *testing.T) {
	// Setup
	mockRepo := new(mocks.BookingRepository)
	mockUserClient := new(usermocks.UserServiceClient)
	mockRideClient := new(ridemocks.RideServiceClient)

	bookingServer := NewBookingServer(mockRepo, mockUserClient, mockRideClient)

	ctx := context.Background()
	req := &pb.CreateBookingRequest{
		UserId: 1,
		Ride: &pb.Ride{
			Source:      "New York",
			Destination: "Boston",
			Distance:    200,
			Cost:        150,
		},
	}

	// Expectations
	mockUserClient.On("GetUser", ctx, &userpb.GetUserRequest{UserId: 1}).
		Return(&userpb.GetUserResponse{Name: "John Doe"}, nil)

	mockRideClient.On("CreateRide", ctx, &ridepb.CreateRideRequest{
		Source:      "New York",
		Destination: "Boston",
		Distance:    200,
		Cost:        150,
	}).Return(nil, errors.New("failed to create ride"))

	// Action
	resp, err := bookingServer.CreateBooking(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verify expectations
	mockUserClient.AssertExpectations(t)
	mockRideClient.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateBooking_RepositoryError(t *testing.T) {
	// Setup
	mockRepo := new(mocks.BookingRepository)
	mockUserClient := new(usermocks.UserServiceClient)
	mockRideClient := new(ridemocks.RideServiceClient)

	bookingServer := NewBookingServer(mockRepo, mockUserClient, mockRideClient)

	ctx := context.Background()
	req := &pb.CreateBookingRequest{
		UserId: 1,
		Ride: &pb.Ride{
			Source:      "New York",
			Destination: "Boston",
			Distance:    200,
			Cost:        150,
		},
	}

	// Expectations
	mockUserClient.On("GetUser", ctx, &userpb.GetUserRequest{UserId: 1}).
		Return(&userpb.GetUserResponse{Name: "John Doe"}, nil)

	mockRideClient.On("CreateRide", ctx, &ridepb.CreateRideRequest{
		Source:      "New York",
		Destination: "Boston",
		Distance:    200,
		Cost:        150,
	}).Return(&ridepb.CreateRideResponse{RideId: 5}, nil)

	mockRepo.On("Create", ctx, int32(1), int32(5)).Return(nil, errors.New("database error"))

	// Action
	resp, err := bookingServer.CreateBooking(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verify expectations
	mockUserClient.AssertExpectations(t)
	mockRideClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestGetBooking_Success(t *testing.T) {
	// Setup
	mockRepo := new(mocks.BookingRepository)
	mockUserClient := new(usermocks.UserServiceClient)
	mockRideClient := new(ridemocks.RideServiceClient)

	bookingServer := NewBookingServer(mockRepo, mockUserClient, mockRideClient)

	ctx := context.Background()
	req := &pb.GetBookingRequest{
		BookingId: 1,
	}

	// Mock the booking
	mockBooking := &repository.Booking{
		ID:     1,
		UserID: 2,
		RideID: 3,
		Time:   "2023-01-01T12:00:00Z",
	}

	// Expectations
	mockRepo.On("GetByID", ctx, int32(1)).Return(mockBooking, nil)

	mockUserClient.On("GetUser", ctx, &userpb.GetUserRequest{UserId: 2}).
		Return(&userpb.GetUserResponse{Name: "John Doe"}, nil)

	mockRideClient.On("GetRide", ctx, &ridepb.GetRideRequest{RideId: 3}).
		Return(&ridepb.Ride{
			RideId:      3,
			Source:      "New York",
			Destination: "Boston",
			Distance:    200,
			Cost:        150,
		}, nil)

	// Action
	resp, err := bookingServer.GetBooking(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "John Doe", resp.Name)
	assert.Equal(t, "New York", resp.Source)
	assert.Equal(t, "Boston", resp.Destination)
	assert.Equal(t, int32(200), resp.Distance)
	assert.Equal(t, int32(150), resp.Cost)
	assert.Equal(t, "2023-01-01T12:00:00Z", resp.Time)

	// Verify expectations
	mockRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockRideClient.AssertExpectations(t)
}

func TestGetBooking_InvalidId(t *testing.T) {
	// Setup
	mockRepo := new(mocks.BookingRepository)
	mockUserClient := new(usermocks.UserServiceClient)
	mockRideClient := new(ridemocks.RideServiceClient)

	bookingServer := NewBookingServer(mockRepo, mockUserClient, mockRideClient)

	ctx := context.Background()
	req := &pb.GetBookingRequest{
		BookingId: 0,
	}

	// Action
	resp, err := bookingServer.GetBooking(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)

	// No repository calls should be made
	mockRepo.AssertNotCalled(t, "GetByID")
	mockUserClient.AssertNotCalled(t, "GetUser")
	mockRideClient.AssertNotCalled(t, "GetRide")
}

func TestGetBooking_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(mocks.BookingRepository)
	mockUserClient := new(usermocks.UserServiceClient)
	mockRideClient := new(ridemocks.RideServiceClient)

	bookingServer := NewBookingServer(mockRepo, mockUserClient, mockRideClient)

	ctx := context.Background()
	req := &pb.GetBookingRequest{
		BookingId: 999,
	}

	// Expectations
	mockRepo.On("GetByID", ctx, int32(999)).Return(nil, errors.New("booking not found"))

	// Action
	resp, err := bookingServer.GetBooking(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verify only repository was called but not the client services
	mockRepo.AssertExpectations(t)
	mockUserClient.AssertNotCalled(t, "GetUser")
	mockRideClient.AssertNotCalled(t, "GetRide")
}

func TestGetBooking_UserServiceError(t *testing.T) {
	// Setup
	mockRepo := new(mocks.BookingRepository)
	mockUserClient := new(usermocks.UserServiceClient)
	mockRideClient := new(ridemocks.RideServiceClient)

	bookingServer := NewBookingServer(mockRepo, mockUserClient, mockRideClient)

	ctx := context.Background()
	req := &pb.GetBookingRequest{
		BookingId: 1,
	}

	// Mock the booking
	mockBooking := &repository.Booking{
		ID:     1,
		UserID: 2,
		RideID: 3,
		Time:   "2023-01-01T12:00:00Z",
	}

	// Expectations
	mockRepo.On("GetByID", ctx, int32(1)).Return(mockBooking, nil)

	mockUserClient.On("GetUser", ctx, &userpb.GetUserRequest{UserId: 2}).
		Return(nil, errors.New("user service error"))

	// Action
	resp, err := bookingServer.GetBooking(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verify expectations - ride service should not be called after user service fails
	mockRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockRideClient.AssertNotCalled(t, "GetRide")
}

func TestGetBooking_RideServiceError(t *testing.T) {
	// Setup
	mockRepo := new(mocks.BookingRepository)
	mockUserClient := new(usermocks.UserServiceClient)
	mockRideClient := new(ridemocks.RideServiceClient)

	bookingServer := NewBookingServer(mockRepo, mockUserClient, mockRideClient)

	ctx := context.Background()
	req := &pb.GetBookingRequest{
		BookingId: 1,
	}

	// Mock the booking
	mockBooking := &repository.Booking{
		ID:     1,
		UserID: 2,
		RideID: 3,
		Time:   "2023-01-01T12:00:00Z",
	}

	// Expectations
	mockRepo.On("GetByID", ctx, int32(1)).Return(mockBooking, nil)

	mockUserClient.On("GetUser", ctx, &userpb.GetUserRequest{UserId: 2}).
		Return(&userpb.GetUserResponse{Name: "John Doe"}, nil)

	mockRideClient.On("GetRide", ctx, &ridepb.GetRideRequest{RideId: 3}).
		Return(nil, errors.New("ride service error"))

	// Action
	resp, err := bookingServer.GetBooking(ctx, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verify expectations
	mockRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockRideClient.AssertExpectations(t)
}
