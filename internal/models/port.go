package models

import (
	"fleet-management/internal/database"
	"time"
)

type Port struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	LocationLat       float64 `json:"location_lat"`
	LocationLng       float64 `json:"location_lng"`
	LocationElevation float64 `json:"location_elevation"`
	TimeCreated       int64   `json:"time_created"`
	TimeUpdated       int64   `json:"time_updated"`
}

func CreatePort(db *database.Database, name string, lat, lng, elevation float64) (*Port, error) {
	now := time.Now().Unix()
	portID := generateID()

	query := `
		INSERT INTO ports (id, name, location_lat, location_lng, location_elevation, time_created, time_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id, name, location_lat, location_lng, location_elevation, time_created, time_updated
	`

	var port Port
	err := db.QueryRow(query, portID, name, lat, lng, elevation, now).Scan(
		&port.ID, &port.Name, &port.LocationLat, &port.LocationLng, &port.LocationElevation, &port.TimeCreated, &port.TimeUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &port, nil
}

func GetPort(db *database.Database, portID string) (*Port, error) {
	query := `
		SELECT id, name, location_lat, location_lng, location_elevation, time_created, time_updated
		FROM ports WHERE id = $1
	`

	var port Port
	err := db.QueryRow(query, portID).Scan(
		&port.ID, &port.Name, &port.LocationLat, &port.LocationLng, &port.LocationElevation, &port.TimeCreated, &port.TimeUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &port, nil
}

func GetPortFleets(db *database.Database, portID string) ([]string, error) {
	query := `SELECT id FROM fleets WHERE port_id = $1`

	rows, err := db.Query(query, portID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fleetIDs []string
	for rows.Next() {
		var fleetID string
		if err := rows.Scan(&fleetID); err != nil {
			return nil, err
		}
		fleetIDs = append(fleetIDs, fleetID)
	}

	return fleetIDs, nil
}
