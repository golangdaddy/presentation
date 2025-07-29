package handlers

import (
	"fleet-management/internal/database"
	"fleet-management/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePort(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name              string  `json:"name" binding:"required"`
			LocationLat       float64 `json:"location_lat" binding:"required"`
			LocationLng       float64 `json:"location_lng" binding:"required"`
			LocationElevation float64 `json:"location_elevation" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		port, err := models.CreatePort(db, req.Name, req.LocationLat, req.LocationLng, req.LocationElevation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create port"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":   port.ID,
			"type": "port",
			"fields": gin.H{
				"name":               port.Name,
				"location_lat":       port.LocationLat,
				"location_lng":       port.LocationLng,
				"location_elevation": port.LocationElevation,
			},
		})
	}
}

func GetPort(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		portID := c.Param("port_id")
		if portID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Port ID is required"})
			return
		}

		port, err := models.GetPort(db, portID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Port not found"})
			return
		}

		// Get associated fleets
		fleets, err := models.GetPortFleets(db, portID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get fleets"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":   port.ID,
			"type": "port",
			"fields": gin.H{
				"name":               port.Name,
				"location_lat":       port.LocationLat,
				"location_lng":       port.LocationLng,
				"location_elevation": port.LocationElevation,
			},
			"fleet": fleets,
		})
	}
}

func AddFleetToPort(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		portID := c.Param("port_id")
		fleetID := c.Param("fleet_id")

		if portID == "" || fleetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Port ID and Fleet ID are required"})
			return
		}

		// Update fleet's port_id
		query := `UPDATE fleets SET port_id = $1 WHERE id = $2`
		_, err := db.Exec(query, portID, fleetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add fleet to port"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Fleet added to port"})
	}
}

func RemoveFleetFromPort(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		portID := c.Param("port_id")
		fleetID := c.Param("fleet_id")

		if portID == "" || fleetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Port ID and Fleet ID are required"})
			return
		}

		// Remove fleet from port
		query := `UPDATE fleets SET port_id = NULL WHERE id = $1 AND port_id = $2`
		_, err := db.Exec(query, fleetID, portID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove fleet from port"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Fleet removed from port"})
	}
}
