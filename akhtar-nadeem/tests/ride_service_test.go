package tests

import (
	"context"
	"errors"
	"github.com/akhtarCareem/golang-assignment/internal/services"
	proto "github.com/akhtarCareem/golang-assignment/proto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRidesRepo struct {
	rides map[string]mockRideData
	err   error
}

type mockRideData struct {
	RideID      string
	Source      string
	Destination string
	Distance    int32
	Cost        int32
}

func (m *mockRidesRepo) UpdateRide(rideID, source, destination string, distance, cost int32) error {
	if m.err != nil {
		return m.err
	}
	ride, ok := m.rides[rideID]
	if !ok {
		return errors.New("ride not found")
	}
	ride.Source = source
	ride.Destination = destination
	ride.Distance = distance
	ride.Cost = cost
	m.rides[rideID] = ride
	return nil
}

func TestRidesService(t *testing.T) {
	repo := &mockRidesRepo{
		rides: map[string]mockRideData{
			"mock_ride_id": {RideID: "mock_ride_id", Source: "X", Destination: "Y", Distance: 5, Cost: 50},
		},
	}
	log := logrus.New()
	svc := services.NewRidesService(repo, log)

	// UpdateRide success
	resp, err := svc.UpdateRide(context.Background(), &proto.UpdateRideRequest{
		RideId: "mock_ride_id",
		Ride: &proto.Ride{
			Source:      "A",
			Destination: "B",
			Distance:    10,
			Cost:        100,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "ride updated successfully", resp.Message)
	assert.Equal(t, "A", repo.rides["mock_ride_id"].Source)
	assert.Equal(t, "B", repo.rides["mock_ride_id"].Destination)
	assert.Equal(t, int32(10), repo.rides["mock_ride_id"].Distance)
	assert.Equal(t, int32(100), repo.rides["mock_ride_id"].Cost)

	// UpdateRide error scenario
	repo.err = errors.New("db error")
	resp, err = svc.UpdateRide(context.Background(), &proto.UpdateRideRequest{
		RideId: "mock_ride_id",
		Ride: &proto.Ride{
			Source:      "C",
			Destination: "D",
			Distance:    20,
			Cost:        200,
		},
	})
	assert.Error(t, err)
	assert.Nil(t, resp)
}
