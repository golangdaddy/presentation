package handlers

import (
	"fleet-management/internal/config"
	"fleet-management/internal/database"
	"fleet-management/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateSession(db *database.Database, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// In a real implementation, you would send a magic link via email
		// For now, we'll just return a success response
		c.JSON(http.StatusOK, gin.H{"message": "Magic link sent to email"})
	}
}

func VerifySession(db *database.Database, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
			return
		}

		// In a real implementation, you would verify the token
		// For now, we'll create a session directly
		now := time.Now().Unix()
		sessionID := time.Now().Format("20060102150405") + "000"
		expiry := now + (24 * 60 * 60) // 24 hours

		// Create a temporary user for demo purposes
		user, err := models.CreateUser(db, "demo@example.com", "viewer")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Create session
		query := `
			INSERT INTO sessions (id, user_id, expiry, time_created, time_updated)
			VALUES ($1, $2, $3, $4, $4)
		`
		_, err = db.Exec(query, sessionID, user.ID, expiry, now)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":      sessionID,
			"user_id": user.ID,
			"expiry":  expiry,
		})
	}
}

func DeleteSession(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("session_id")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
			return
		}

		query := `DELETE FROM sessions WHERE id = $1`
		_, err := db.Exec(query, sessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Session deleted"})
	}
}
