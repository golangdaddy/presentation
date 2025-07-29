package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Response struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})

	// Auth endpoints
	http.HandleFunc("/api/v1/auth/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"Magic link sent to email"}`)
	})

	http.HandleFunc("/api/v1/auth/session/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, `{"error":"Token is required"}`, http.StatusBadRequest)
			return
		}

		// Create a demo session
		sessionID := time.Now().Format("20060102150405") + "000"
		expiry := time.Now().Add(24 * time.Hour).Unix()

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"id":"%s",
			"user_id":"demo_user_123",
			"expiry":%d
		}`, sessionID, expiry)
	})

	// User management endpoints
	http.HandleFunc("/api/v1/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var req struct {
				Email string `json:"email"`
				Role  string `json:"role"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
				return
			}

			userID := time.Now().Format("20060102150405") + "000"
			response := Response{
				ID:   userID,
				Type: "user",
				Fields: map[string]interface{}{
					"email": req.Email,
					"role":  req.Role,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Session management
	http.HandleFunc("/api/v1/auth/session/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 5 {
			http.Error(w, `{"error":"Session ID required"}`, http.StatusBadRequest)
			return
		}

		_ = pathParts[4] // sessionID
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"Session deleted"}`)
	})

	// Port management endpoints
	http.HandleFunc("/api/v1/port", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Name              string  `json:"name"`
			LocationLat       float64 `json:"location_lat"`
			LocationLng       float64 `json:"location_lng"`
			LocationElevation float64 `json:"location_elevation"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		portID := time.Now().Format("20060102150405") + "000"
		response := Response{
			ID:   portID,
			Type: "port",
			Fields: map[string]interface{}{
				"name":               req.Name,
				"location_lat":       req.LocationLat,
				"location_lng":       req.LocationLng,
				"location_elevation": req.LocationElevation,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Port detail and fleet management
	http.HandleFunc("/api/v1/port/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, `{"error":"Port ID required"}`, http.StatusBadRequest)
			return
		}

		portID := pathParts[3]
		if strings.Contains(portID, "/") {
			portID = strings.Split(portID, "/")[0]
		}

		if r.Method == "GET" {
			// Get port details
			response := Response{
				ID:   portID,
				Type: "port",
				Fields: map[string]interface{}{
					"name":               "Demo Port",
					"location_lat":       51.3308,
					"location_lng":       0.0323,
					"location_elevation": 183.0,
				},
			}

			fullResponse := map[string]interface{}{
				"id":     response.ID,
				"type":   response.Type,
				"fields": response.Fields,
				"fleet":  []string{"fleet_123", "fleet_456"},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(fullResponse)
		} else if r.Method == "PUT" && len(pathParts) >= 6 && pathParts[5] == "fleet" {
			// Add fleet to port
			fleetID := pathParts[6]
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"message":"Fleet %s added to port %s"}`, fleetID, portID)
		} else if r.Method == "DELETE" && len(pathParts) >= 6 && pathParts[5] == "fleet" {
			// Remove fleet from port
			fleetID := pathParts[6]
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"message":"Fleet %s removed from port %s"}`, fleetID, portID)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Asset template endpoints
	http.HandleFunc("/api/v1/asset-template", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Name           string `json:"name"`
			ManufacturerID string `json:"manufacturer_id"`
			ProductWidth   int    `json:"product_width"`
			ProductHeight  int    `json:"product_height"`
			ProductLength  int    `json:"product_length"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		templateID := time.Now().Format("20060102150405") + "000"
		response := Response{
			ID:   templateID,
			Type: "template",
			Fields: map[string]interface{}{
				"name":            req.Name,
				"manufacturer_id": req.ManufacturerID,
				"product_width":   req.ProductWidth,
				"product_height":  req.ProductHeight,
				"product_length":  req.ProductLength,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Asset template detail and component management
	http.HandleFunc("/api/v1/asset-template/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, `{"error":"Template ID required"}`, http.StatusBadRequest)
			return
		}

		templateID := pathParts[3]
		if strings.Contains(templateID, "/") {
			templateID = strings.Split(templateID, "/")[0]
		}

		if r.Method == "GET" {
			// Get template details
			response := Response{
				ID:   templateID,
				Type: "template",
				Fields: map[string]interface{}{
					"name":            "HyperDrone XF-11",
					"manufacturer_id": "demo_manufacturer",
					"product_width":   80,
					"product_height":  40,
					"product_length":  60,
				},
			}

			fullResponse := map[string]interface{}{
				"id":         response.ID,
				"type":       response.Type,
				"fields":     response.Fields,
				"components": []string{"component_123", "component_456"},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(fullResponse)
		} else if r.Method == "PATCH" {
			// Update template
			var req struct {
				ProductWeight int `json:"product_weight"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"message":"Template %s updated"}`, templateID)
		} else if len(pathParts) >= 6 && pathParts[5] == "component" {
			if r.Method == "POST" {
				// Create component
				var req struct {
					Name           string `json:"name"`
					ManufacturerID string `json:"manufacturer_id"`
					ProductWeight  int    `json:"product_weight"`
				}

				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
					return
				}

				componentID := time.Now().Format("20060102150405") + "000"
				response := Response{
					ID:   componentID,
					Type: "component",
					Fields: map[string]interface{}{
						"name":            req.Name,
						"manufacturer_id": req.ManufacturerID,
						"product_weight":  req.ProductWeight,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else if r.Method == "PATCH" && len(pathParts) >= 7 {
				// Update component
				componentID := pathParts[6]
				var req struct {
					ProductWeight int `json:"product_weight"`
				}

				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"message":"Component %s updated"}`, componentID)
			} else if r.Method == "DELETE" && len(pathParts) >= 7 {
				// Delete component
				componentID := pathParts[6]
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"message":"Component %s deleted"}`, componentID)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Fleet management endpoints
	http.HandleFunc("/api/v1/fleet", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Create fleet
			var req struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
				return
			}

			if req.Name == "" {
				http.Error(w, `{"error":"Name is required"}`, http.StatusBadRequest)
				return
			}

			fleetID := time.Now().Format("20060102150405") + "000"
			response := Response{
				ID:   fleetID,
				Type: "fleet",
				Fields: map[string]interface{}{
					"name":        req.Name,
					"description": req.Description,
					"port":        "demo_port_123",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Fleet detail and management
	http.HandleFunc("/api/v1/fleet/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, `{"error":"Fleet ID required"}`, http.StatusBadRequest)
			return
		}

		fleetID := pathParts[3]
		if strings.Contains(fleetID, "/") {
			fleetID = strings.Split(fleetID, "/")[0]
		}

		// Check if this is a compliance request
		if len(pathParts) >= 5 && pathParts[4] == "compliance" {
			compliance := map[string]interface{}{
				"fleet_id": fleetID,
				"assets_with_overdue_parts": []map[string]interface{}{
					{
						"asset_id":   "demo_asset_123",
						"asset_name": "Demo Drone A",
						"overdue_parts": []map[string]interface{}{
							{
								"part_id":              "demo_part_123",
								"part_name":            "Front-left rotor blades",
								"last_inspection_time": 1690867200000,
								"inspection_frequency": 30,
								"days_overdue":         15,
							},
						},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(compliance)
			return
		}

		// Check for template management
		if len(pathParts) >= 6 && pathParts[5] == "template" {
			templateID := pathParts[6]
			if r.Method == "PUT" {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"message":"Template %s added to fleet %s"}`, templateID, fleetID)
			} else if r.Method == "DELETE" {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"message":"Template %s removed from fleet %s"}`, templateID, fleetID)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Check for asset management
		if len(pathParts) >= 6 && pathParts[5] == "asset" {
			assetID := pathParts[6]
			if r.Method == "PUT" {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"message":"Asset %s added to fleet %s"}`, assetID, fleetID)
			} else if r.Method == "DELETE" {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"message":"Asset %s removed from fleet %s"}`, assetID, fleetID)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Regular fleet detail response
		response := Response{
			ID:   fleetID,
			Type: "fleet",
			Fields: map[string]interface{}{
				"name":        "Demo Fleet",
				"description": "A demo fleet for testing",
				"port":        "demo_port_123",
			},
		}

		// Add templates and assets arrays
		fullResponse := map[string]interface{}{
			"id":        response.ID,
			"type":      response.Type,
			"fields":    response.Fields,
			"templates": []string{"template_123", "template_456"},
			"assets":    []string{"asset_123", "asset_456"},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fullResponse)
	})

	// Asset endpoints
	http.HandleFunc("/api/v1/asset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Name     string `json:"name"`
			Template string `json:"template"`
			DateBuy  int64  `json:"date_buy"`
			Warranty string `json:"warranty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		assetID := time.Now().Format("20060102150405") + "000"
		response := Response{
			ID:   assetID,
			Type: "asset",
			Fields: map[string]interface{}{
				"name":     req.Name,
				"template": req.Template,
				"date_buy": req.DateBuy,
				"warranty": req.Warranty,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Asset detail and management
	http.HandleFunc("/api/v1/asset/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, `{"error":"Asset ID required"}`, http.StatusBadRequest)
			return
		}

		assetID := pathParts[3]
		if strings.Contains(assetID, "/") {
			assetID = strings.Split(assetID, "/")[0]
		}

		if r.Method == "GET" {
			// Get asset details with components
			response := Response{
				ID:   assetID,
				Type: "asset",
				Fields: map[string]interface{}{
					"name":     "My New Drone A",
					"template": "hyperdrone_xf-11",
					"date_buy": 73287028340,
					"warranty": "1y",
				},
			}

			components := []map[string]interface{}{
				{
					"id":                   "part_123",
					"name":                 "Front-left rotor blades",
					"serial_number":        "12345",
					"condition":            "new",
					"notes":                "Installed",
					"inspection_frequency": 30,
					"attachments":          []string{"https://storage.example.com/attachments/part_photo1.jpg"},
				},
				{
					"id":                   "part_456",
					"name":                 "Battery pack",
					"serial_number":        "67890",
					"condition":            "functional",
					"notes":                "Fully charged",
					"inspection_frequency": 14,
					"attachments":          []string{},
				},
			}

			fullResponse := map[string]interface{}{
				"id":          response.ID,
				"type":        response.Type,
				"fields":      response.Fields,
				"attachments": []string{"https://storage.example.com/attachments/asset_photo1.jpg"},
				"components":  components,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(fullResponse)
		} else if len(pathParts) >= 6 && pathParts[5] == "attachment" {
			if r.Method == "POST" {
				// Add attachment
				var req struct {
					URI      string `json:"uri"`
					Name     string `json:"name"`
					MimeType string `json:"mime_type"`
				}

				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
					return
				}

				attachmentID := time.Now().Format("20060102150405") + "000"
				response := map[string]interface{}{
					"id":  attachmentID,
					"uri": req.URI,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else if r.Method == "DELETE" && len(pathParts) >= 7 {
				// Delete attachment
				attachmentID := pathParts[6]
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"message":"Attachment %s deleted"}`, attachmentID)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else if len(pathParts) >= 6 && pathParts[5] == "inspections" {
			if len(pathParts) >= 7 && pathParts[6] == "schedule" && r.Method == "POST" {
				// Schedule inspection
				var req struct {
					Timestamp int64 `json:"timestamp"`
				}

				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
					return
				}

				inspectionID := time.Now().Format("20060102150405") + "000"
				response := Response{
					ID:   inspectionID,
					Type: "inspection",
					Fields: map[string]interface{}{
						"asset_id": assetID,
						"time":     req.Timestamp,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else if len(pathParts) >= 7 && pathParts[6] == "log" && r.Method == "POST" {
				// Log inspection
				var req []map[string]interface{}

				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
					return
				}

				var responses []map[string]interface{}
				for _, inspection := range req {
					inspectionID := time.Now().Format("20060102150405") + "000"
					response := map[string]interface{}{
						"id":     inspectionID,
						"type":   "inspection",
						"fields": inspection,
					}
					responses = append(responses, response)
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(responses)
			} else if r.Method == "GET" {
				// Get inspections
				inspections := []map[string]interface{}{
					{
						"id": "inspection_123",
						"fields": map[string]interface{}{
							"asset_id": assetID,
							"time":     7978894354358,
						},
						"attachments": []string{"https://storage.example.com/attachments/inspection_photo1.jpg"},
					},
					{
						"id": "inspection_456",
						"fields": map[string]interface{}{
							"asset_id":      assetID,
							"component_id":  "component_123",
							"serial_number": "656172e88398d4",
							"action":        "replaced",
							"condition":     "new",
							"notes":         "was damaged",
							"time":          78269786307,
						},
						"attachments": []string{
							"https://storage.example.com/attachments/replacement_evidence.pdf",
							"https://storage.example.com/attachments/damage_photo.jpg",
						},
					},
				}

				response := map[string]interface{}{
					"id":          assetID,
					"type":        "asset",
					"inspections": inspections,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Inspection endpoints
	http.HandleFunc("/api/v1/inspection/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, `{"error":"Inspection ID required"}`, http.StatusBadRequest)
			return
		}

		inspectionID := pathParts[3]
		if strings.Contains(inspectionID, "/") {
			inspectionID = strings.Split(inspectionID, "/")[0]
		}

		if r.Method == "PATCH" {
			// Update inspection
			var req struct {
				Condition string `json:"condition"`
				Notes     string `json:"notes"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"message":"Inspection %s updated"}`, inspectionID)
		} else if r.Method == "DELETE" {
			// Delete inspection
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"message":"Inspection %s deleted"}`, inspectionID)
		} else if len(pathParts) >= 6 && pathParts[5] == "attachment" {
			if r.Method == "POST" {
				// Add attachment to inspection
				var req struct {
					URI      string `json:"uri"`
					Name     string `json:"name"`
					MimeType string `json:"mime_type"`
				}

				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
					return
				}

				attachmentID := time.Now().Format("20060102150405") + "000"
				response := map[string]interface{}{
					"id":  attachmentID,
					"uri": req.URI,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else if r.Method == "DELETE" && len(pathParts) >= 7 {
				// Delete attachment from inspection
				attachmentID := pathParts[6]
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"message":"Attachment %s deleted from inspection %s"}`, attachmentID, inspectionID)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Asset part endpoints
	http.HandleFunc("/api/v1/asset-part/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, `{"error":"Asset Part ID required"}`, http.StatusBadRequest)
			return
		}

		partID := pathParts[3]
		if strings.Contains(partID, "/") {
			partID = strings.Split(partID, "/")[0]
		}

		if r.Method == "GET" {
			// Get asset part details
			response := Response{
				ID:   partID,
				Type: "asset-part",
				Fields: map[string]interface{}{
					"name":                 "Front-left rotor blades",
					"serial_number":        "12345",
					"condition":            "new",
					"notes":                "Installed",
					"inspection_frequency": 30,
				},
			}

			fullResponse := map[string]interface{}{
				"id":          response.ID,
				"type":        response.Type,
				"fields":      response.Fields,
				"attachments": []string{"https://storage.example.com/attachments/part_photo1.jpg"},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(fullResponse)
		} else if r.Method == "PATCH" {
			// Update asset part
			var req struct {
				InspectionFrequency int64 `json:"inspection_frequency"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"message":"Asset part %s updated"}`, partID)
		} else if len(pathParts) >= 6 && pathParts[5] == "attachment" {
			if r.Method == "POST" {
				// Add attachment to asset part
				var req struct {
					URI      string `json:"uri"`
					Name     string `json:"name"`
					MimeType string `json:"mime_type"`
				}

				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
					return
				}

				attachmentID := time.Now().Format("20060102150405") + "000"
				response := map[string]interface{}{
					"id":  attachmentID,
					"uri": req.URI,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else if r.Method == "DELETE" && len(pathParts) >= 7 {
				// Delete attachment from asset part
				attachmentID := pathParts[6]
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"message":"Attachment %s deleted from asset part %s"}`, attachmentID, partID)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// User detail endpoints
	http.HandleFunc("/api/v1/user/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, `{"error":"User ID required"}`, http.StatusBadRequest)
			return
		}

		userID := pathParts[3]
		if strings.Contains(userID, "/") {
			userID = strings.Split(userID, "/")[0]
		}

		if r.Method == "GET" {
			// Get user details
			response := Response{
				ID:   userID,
				Type: "user",
				Fields: map[string]interface{}{
					"email": "demo@example.com",
					"role":  "reporter",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else if r.Method == "PATCH" {
			// Update user
			var req struct {
				Role string `json:"role"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"message":"User %s updated"}`, userID)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
