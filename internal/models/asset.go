package models

import (
	"fleet-management/internal/database"
	"time"
)

type Asset struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	TemplateID  string  `json:"template_id"`
	FleetID     *string `json:"fleet_id"`
	DateBuy     *int64  `json:"date_buy"`
	DateInstall *int64  `json:"date_install"`
	Warranty    *string `json:"warranty"`
}

type AssetPart struct {
	ID                  string  `json:"id"`
	AssetID             string  `json:"asset_id"`
	ComponentID         string  `json:"component_id"`
	Name                string  `json:"name"`
	SerialNumber        *string `json:"serial_number"`
	Condition           *string `json:"condition"`
	Notes               *string `json:"notes"`
	InspectionFrequency *int64  `json:"inspection_frequency"`
}

type AssetTemplate struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ManufacturerID string `json:"manufacturer_id"`
	ProductWeight  *int64 `json:"product_weight"`
	ProductWidth   *int   `json:"product_width"`
	ProductHeight  *int   `json:"product_height"`
	ProductLength  *int   `json:"product_length"`
	TimeCreated    int64  `json:"time_created"`
	TimeUpdated    int64  `json:"time_updated"`
}

func CreateAsset(db *database.Database, name, templateID string, fleetID *string, dateBuy, dateInstall *int64, warranty *string) (*Asset, error) {
	assetID := generateID()

	query := `
		INSERT INTO assets (id, name, template_id, fleet_id, date_buy, date_install, warranty)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, template_id, fleet_id, date_buy, date_install, warranty
	`

	var asset Asset
	err := db.QueryRow(query, assetID, name, templateID, fleetID, dateBuy, dateInstall, warranty).Scan(
		&asset.ID, &asset.Name, &asset.TemplateID, &asset.FleetID, &asset.DateBuy, &asset.DateInstall, &asset.Warranty,
	)

	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func GetAsset(db *database.Database, assetID string) (*Asset, error) {
	query := `
		SELECT id, name, template_id, fleet_id, date_buy, date_install, warranty
		FROM assets WHERE id = $1
	`

	var asset Asset
	err := db.QueryRow(query, assetID).Scan(
		&asset.ID, &asset.Name, &asset.TemplateID, &asset.FleetID, &asset.DateBuy, &asset.DateInstall, &asset.Warranty,
	)

	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func CreateAssetTemplate(db *database.Database, name, manufacturerID string, productWeight *int64, productWidth, productHeight, productLength *int) (*AssetTemplate, error) {
	now := time.Now().Unix()
	templateID := generateID()

	query := `
		INSERT INTO assettemplates (id, name, manufacturer_id, product_weight, product_width, product_height, product_length, time_created, time_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
		RETURNING id, name, manufacturer_id, product_weight, product_width, product_height, product_length, time_created, time_updated
	`

	var template AssetTemplate
	err := db.QueryRow(query, templateID, name, manufacturerID, productWeight, productWidth, productHeight, productLength, now).Scan(
		&template.ID, &template.Name, &template.ManufacturerID, &template.ProductWeight, &template.ProductWidth, &template.ProductHeight, &template.ProductLength, &template.TimeCreated, &template.TimeUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &template, nil
}

func GetAssetTemplate(db *database.Database, templateID string) (*AssetTemplate, error) {
	query := `
		SELECT id, name, manufacturer_id, product_weight, product_width, product_height, product_length, time_created, time_updated
		FROM assettemplates WHERE id = $1
	`

	var template AssetTemplate
	err := db.QueryRow(query, templateID).Scan(
		&template.ID, &template.Name, &template.ManufacturerID, &template.ProductWeight, &template.ProductWidth, &template.ProductHeight, &template.ProductLength, &template.TimeCreated, &template.TimeUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &template, nil
}

func GetAssetAttachments(db *database.Database, assetID string) ([]string, error) {
	query := `SELECT uri FROM attachments WHERE entity_type = 'Asset' AND entity_id = $1`

	rows, err := db.Query(query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uris []string
	for rows.Next() {
		var uri string
		if err := rows.Scan(&uri); err != nil {
			return nil, err
		}
		uris = append(uris, uri)
	}

	return uris, nil
}

func GetAssetParts(db *database.Database, assetID string) ([]AssetPart, error) {
	query := `
		SELECT id, asset_id, component_id, name, serial_number, condition, notes, inspection_frequency
		FROM assetparts WHERE asset_id = $1
	`

	rows, err := db.Query(query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []AssetPart
	for rows.Next() {
		var part AssetPart
		err := rows.Scan(
			&part.ID, &part.AssetID, &part.ComponentID, &part.Name, &part.SerialNumber, &part.Condition, &part.Notes, &part.InspectionFrequency,
		)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	return parts, nil
}
