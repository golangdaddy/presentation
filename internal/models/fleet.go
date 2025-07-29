package models

import (
	"fleet-management/internal/database"
	"time"
)

type Fleet struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PortID      string `json:"port_id"`
	TimeCreated int64  `json:"time_created"`
	TimeUpdated int64  `json:"time_updated"`
}

func CreateFleet(db *database.Database, name, description, portID string) (*Fleet, error) {
	now := time.Now().Unix()
	fleetID := generateID()

	query := `
		INSERT INTO fleets (id, name, description, port_id, time_created, time_updated)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id, name, description, port_id, time_created, time_updated
	`

	var fleet Fleet
	err := db.QueryRow(query, fleetID, name, description, portID, now).Scan(
		&fleet.ID, &fleet.Name, &fleet.Description, &fleet.PortID, &fleet.TimeCreated, &fleet.TimeUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &fleet, nil
}

func GetFleet(db *database.Database, fleetID string) (*Fleet, error) {
	query := `
		SELECT id, name, description, port_id, time_created, time_updated
		FROM fleets WHERE id = $1
	`

	var fleet Fleet
	err := db.QueryRow(query, fleetID).Scan(
		&fleet.ID, &fleet.Name, &fleet.Description, &fleet.PortID, &fleet.TimeCreated, &fleet.TimeUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &fleet, nil
}

func GetFleetTemplates(db *database.Database, fleetID string) ([]string, error) {
	query := `SELECT template_id FROM fleetassettemplates WHERE fleet_id = $1`

	rows, err := db.Query(query, fleetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templateIDs []string
	for rows.Next() {
		var templateID string
		if err := rows.Scan(&templateID); err != nil {
			return nil, err
		}
		templateIDs = append(templateIDs, templateID)
	}

	return templateIDs, nil
}

func GetFleetAssets(db *database.Database, fleetID string) ([]string, error) {
	query := `SELECT id FROM assets WHERE fleet_id = $1`

	rows, err := db.Query(query, fleetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assetIDs []string
	for rows.Next() {
		var assetID string
		if err := rows.Scan(&assetID); err != nil {
			return nil, err
		}
		assetIDs = append(assetIDs, assetID)
	}

	return assetIDs, nil
}

func AddTemplateToFleet(db *database.Database, fleetID, templateID string) error {
	query := `INSERT INTO fleetassettemplates (fleet_id, template_id) VALUES ($1, $2)`
	_, err := db.Exec(query, fleetID, templateID)
	return err
}

func RemoveTemplateFromFleet(db *database.Database, fleetID, templateID string) error {
	query := `DELETE FROM fleetassettemplates WHERE fleet_id = $1 AND template_id = $2`
	_, err := db.Exec(query, fleetID, templateID)
	return err
}

func AddAssetToFleet(db *database.Database, fleetID, assetID string) error {
	query := `UPDATE assets SET fleet_id = $1 WHERE id = $2`
	_, err := db.Exec(query, fleetID, assetID)
	return err
}

func RemoveAssetFromFleet(db *database.Database, fleetID, assetID string) error {
	query := `UPDATE assets SET fleet_id = NULL WHERE id = $1 AND fleet_id = $2`
	_, err := db.Exec(query, assetID, fleetID)
	return err
}
