package handlers

import (
	"fleet-management/internal/config"
	"fleet-management/internal/database"
	"fleet-management/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(db *database.Database, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
			Role  string `json:"role" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate role
		validRoles := map[string]bool{"viewer": true, "reporter": true, "editor": true, "owner": true}
		if !validRoles[req.Role] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be viewer, reporter, editor, or owner"})
			return
		}

		user, err := models.CreateUser(db, req.Email, req.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":   user.ID,
			"type": "user",
			"fields": gin.H{
				"email": user.Email,
				"role":  user.RoleID,
			},
		})
	}
}

func GetUser(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		user, err := models.GetUser(db, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":   user.ID,
			"type": "user",
			"fields": gin.H{
				"email": user.Email,
				"role":  user.RoleID,
			},
		})
	}
}

func UpdateUser(db *database.Database, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		var req struct {
			Role string `json:"role"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate role if provided
		if req.Role != "" {
			validRoles := map[string]bool{"viewer": true, "reporter": true, "editor": true, "owner": true}
			if !validRoles[req.Role] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be viewer, reporter, editor, or owner"})
				return
			}
		}

		err := models.UpdateUser(db, userID, req.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}
