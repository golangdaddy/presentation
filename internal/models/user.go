package models

import (
	"database/sql"
	"fleet-management/internal/database"
	"time"
)

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	RoleID      string `json:"role_id"`
	TimeCreated int64  `json:"time_created"`
	TimeUpdated int64  `json:"time_updated"`
}

type Session struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Expiry      int64  `json:"expiry"`
	TimeCreated int64  `json:"time_created"`
	TimeUpdated int64  `json:"time_updated"`
}

func VerifySession(db *database.Database, sessionID string) (*User, error) {
	query := `
		SELECT u.id, u.email, u.role_id, u.time_created, u.time_updated
		FROM users u
		JOIN sessions s ON u.id = s.user_id
		WHERE s.id = $1 AND s.expiry > $2
	`

	now := time.Now().Unix()
	var user User
	err := db.QueryRow(query, sessionID, now).Scan(
		&user.ID, &user.Email, &user.RoleID, &user.TimeCreated, &user.TimeUpdated,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func CreateUser(db *database.Database, email, roleID string) (*User, error) {
	now := time.Now().Unix()
	userID := generateID()

	query := `
		INSERT INTO users (id, email, role_id, time_created, time_updated)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, email, role_id, time_created, time_updated
	`

	var user User
	err := db.QueryRow(query, userID, email, roleID, now).Scan(
		&user.ID, &user.Email, &user.RoleID, &user.TimeCreated, &user.TimeUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUser(db *database.Database, userID string) (*User, error) {
	query := `
		SELECT id, email, role_id, time_created, time_updated
		FROM users WHERE id = $1
	`

	var user User
	err := db.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.RoleID, &user.TimeCreated, &user.TimeUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(db *database.Database, userID, roleID string) error {
	now := time.Now().Unix()
	query := `UPDATE users SET role_id = $1, time_updated = $2 WHERE id = $3`
	_, err := db.Exec(query, roleID, now, userID)
	return err
}

func generateID() string {
	// Simple ID generation - in production, use a proper UUID library
	return time.Now().Format("20060102150405") + "000"
}
