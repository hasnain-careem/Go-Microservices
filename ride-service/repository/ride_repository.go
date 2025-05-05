package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Ride struct {
	ID          int32
	Source      string
	Destination string
	Distance    int32
	Cost        int32
}

type RideRepository interface {
	Create(ctx context.Context, source, destination string, distance, cost int32) (int32, error)
	GetByID(ctx context.Context, id int32) (*Ride, error)
	Update(ctx context.Context, id int32, source, destination string, distance, cost int32) (string, error)
}

type PostgresRideRepository struct {
	db *sql.DB
}

func NewPostgresRideRepository(db *sql.DB) RideRepository {
	return &PostgresRideRepository{db: db}
}

func (r *PostgresRideRepository) Create(ctx context.Context, source, destination string, distance, cost int32) (int32, error) {
	query := `INSERT INTO rides (source, destination, distance, cost) VALUES ($1, $2, $3, $4) RETURNING ride_id`
	var rideID int32
	err := r.db.QueryRowContext(ctx, query, source, destination, distance, cost).Scan(&rideID)
	if err != nil {
		log.Printf("Create ride failed: %v", err)
		return 0, err
	}
	return rideID, nil
}

func (r *PostgresRideRepository) GetByID(ctx context.Context, id int32) (*Ride, error) {
	query := `SELECT source, destination, distance, cost FROM rides WHERE ride_id = $1`
	var source, destination string
	var distance, cost int32

	err := r.db.QueryRowContext(ctx, query, id).Scan(&source, &destination, &distance, &cost)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ride not found")
		}
		log.Printf("Get ride failed: %v", err)
		return nil, err
	}

	return &Ride{
		ID:          id,
		Source:      source,
		Destination: destination,
		Distance:    distance,
		Cost:        cost,
	}, nil
}

func (r *PostgresRideRepository) Update(ctx context.Context, id int32, source, destination string, distance, cost int32) (string, error) {
	query := `UPDATE rides SET source = $1, destination = $2, distance = $3, cost = $4 WHERE ride_id = $5`
	_, err := r.db.ExecContext(ctx, query, source, destination, distance, cost, id)
	if err != nil {
		log.Printf("Update ride failed: %v", err)
		return "", err
	}
	return fmt.Sprintf("Ride %d updated successfully", id), nil
}
