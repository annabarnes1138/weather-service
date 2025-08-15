package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"weather-service/pkg/weather"
)

// Server represents the HTTP server
type Server struct {
	weatherService *weather.Service
}

// NewServer creates a new server instance
func NewServer() *Server {
	return &Server{
		weatherService: weather.NewService(),
	}
}

// weatherHandler handles the /weather endpoint
func (s *Server) weatherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse latitude and longitude from query parameters
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	if latStr == "" || lngStr == "" {
		http.Error(w, "Missing required parameters: lat and lng", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude parameter", http.StatusBadRequest)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude parameter", http.StatusBadRequest)
		return
	}

	// Validate coordinate ranges
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		http.Error(w, "Invalid coordinate ranges", http.StatusBadRequest)
		return
	}

	// Get weather forecast
	weatherData, err := s.weatherService.GetForecast(lat, lng)
	if err != nil {
		log.Printf("Error getting weather: %v", err)
		http.Error(w, "Failed to get weather forecast", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weatherData); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// healthHandler provides a health check endpoint
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// rootHandler provides usage instructions
func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, `Weather Service

Usage: GET /weather?lat={latitude}&lng={longitude}

Example: /weather?lat=40.7128&lng=-74.0060

Returns JSON with:
- forecast: Short weather description
- temperature: "hot", "cold", or "moderate"
- location: Coordinates used

Health check: GET /health
`)
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	http.HandleFunc("/weather", s.weatherHandler)
	http.HandleFunc("/health", s.healthHandler)
	http.HandleFunc("/", s.rootHandler)
}

func main() {
	server := NewServer()
	server.setupRoutes()

	port := ":8080"
	log.Printf("Weather service starting on port %s", port)
	log.Printf("Try: http://localhost%s/weather?lat=40.7128&lng=-74.0060", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
