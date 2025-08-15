package weather

import (
	"fmt"
	"strings"
)

// Response represents our API response structure
type Response struct {
	Forecast    string `json:"forecast"`
	Temperature string `json:"temperature"`
	Location    string `json:"location"`
}

// TemperatureCategory represents temperature classification
type TemperatureCategory string

const (
	Hot      TemperatureCategory = "hot"
	Cold     TemperatureCategory = "cold"
	Moderate TemperatureCategory = "moderate"
)

// Temperature thresholds in Fahrenheit
const (
	HotThreshold  = 80
	ColdThreshold = 50
)

// Service provides weather forecast operations
type Service struct {
	client *Client
}

// NewService creates a new weather service
func NewService() *Service {
	return &Service{
		client: NewClient(),
	}
}

// GetForecast retrieves weather forecast for given coordinates
func (s *Service) GetForecast(lat, lng float64) (*Response, error) {
	// Step 1: Get grid coordinates
	gridResp, err := s.client.GetGridCoordinates(lat, lng)
	if err != nil {
		return nil, fmt.Errorf("failed to get grid coordinates: %w", err)
	}

	// Step 2: Get forecast for grid coordinates
	forecastResp, err := s.client.GetForecast(
		gridResp.Properties.GridID,
		gridResp.Properties.GridX,
		gridResp.Properties.GridY,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get forecast: %w", err)
	}

	// Step 3: Find today's forecast
	period, err := s.findTodaysPeriod(forecastResp.Properties.Periods)
	if err != nil {
		return nil, err
	}

	tempCategory := s.categorizeTemperature(period.Temperature)

	return &Response{
		Forecast:    period.ShortForecast,
		Temperature: string(tempCategory),
		Location:    fmt.Sprintf("%.4f, %.4f", lat, lng),
	}, nil
}

// findTodaysPeriod finds the most appropriate period representing "today"
func (s *Service) findTodaysPeriod(periods []Period) (*Period, error) {
	if len(periods) == 0 {
		return nil, fmt.Errorf("no forecast periods found")
	}

	// Get the first period (should be today or tonight)
	period := periods[0]

	// If the first period is tonight, try to get the day period if available
	if !period.IsDaytime && len(periods) > 1 {
		// Check if the next period is for today during daytime
		if strings.Contains(strings.ToLower(periods[1].Name), "today") {
			period = periods[1]
		}
	}

	return &period, nil
}

// categorizeTemperature categorizes temperature as hot, cold, or moderate
func (s *Service) categorizeTemperature(tempF int) TemperatureCategory {
	switch {
	case tempF >= HotThreshold:
		return Hot
	case tempF <= ColdThreshold:
		return Cold
	default:
		return Moderate
	}
}
