package repositories

import (
	"github.com/akhtarCareem/golang-assignment/internal/database"
	"github.com/jinzhu/gorm"
)

type RidesRepository interface {
	UpdateRide(rideID, source, destination string, distance, cost int32) error
}

type ridesRepo struct {
	db *gorm.DB
}

func NewRidesRepository(db *gorm.DB) RidesRepository {
	return &ridesRepo{db: db}
}

func (r *ridesRepo) UpdateRide(rideID, source, destination string, distance, cost int32) error {
	var ride database.Ride
	if err := r.db.Where("ride_id = ?", rideID).First(&ride).Error; err != nil {
		return err
	}
	ride.Source = source
	ride.Destination = destination
	ride.Distance = distance
	ride.Cost = cost
	return r.db.Save(&ride).Error
}
