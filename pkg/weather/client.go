package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Client handles communication with the National Weather Service API
type Client struct {
	baseURL   string
	userAgent string
}

// NewClient creates a new NWS API client
func NewClient() *Client {
	return &Client{
		baseURL:   "https://api.weather.gov",
		userAgent: "WeatherService/1.0 (contact@example.com)",
	}
}

// GridResponse represents the response from NWS points API
type GridResponse struct {
	Properties struct {
		GridID string `json:"gridId"`
		GridX  int    `json:"gridX"`
		GridY  int    `json:"gridY"`
	} `json:"properties"`
}

// ForecastResponse represents the response from NWS forecast API
type ForecastResponse struct {
	Properties struct {
		Periods []Period `json:"periods"`
	} `json:"properties"`
}

// Period represents a forecast period
type Period struct {
	Name          string `json:"name"`
	Temperature   int    `json:"temperature"`
	ShortForecast string `json:"shortForecast"`
	IsDaytime     bool   `json:"isDaytime"`
}

// GetGridCoordinates retrieves grid coordinates for the given lat/lng
func (c *Client) GetGridCoordinates(lat, lng float64) (*GridResponse, error) {
	url := fmt.Sprintf("%s/points/%.4f,%.4f", c.baseURL, lat, lng)

	// Note: Using http.Get for simplicity - in production, would use custom client with timeout
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NWS API returned status %d", resp.StatusCode)
	}

	var gridResp GridResponse
	if err := json.NewDecoder(resp.Body).Decode(&gridResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &gridResp, nil
}

// GetForecast retrieves the forecast for given grid coordinates
func (c *Client) GetForecast(gridID string, gridX, gridY int) (*ForecastResponse, error) {
	url := fmt.Sprintf("%s/gridpoints/%s/%d,%d/forecast", c.baseURL, gridID, gridX, gridY)

	// Note: Using http.Get for simplicity - in production, would use custom client with timeout
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NWS forecast API returned status %d", resp.StatusCode)
	}

	var forecastResp ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &forecastResp, nil
}
