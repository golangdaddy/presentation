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

	// Fleet detail and compliance endpoints
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

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
