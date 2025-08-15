package weather

import (
	"testing"
)

func TestCategorizeTemperature(t *testing.T) {
	service := NewService()

	tests := []struct {
		name     string
		temp     int
		expected TemperatureCategory
	}{
		{"Hot temperature", 85, Hot},
		{"Hot threshold", 80, Hot},
		{"Moderate high", 75, Moderate},
		{"Moderate mid", 65, Moderate},
		{"Moderate low", 51, Moderate},
		{"Cold threshold", 50, Cold},
		{"Cold temperature", 30, Cold},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.categorizeTemperature(tt.temp)
			if result != tt.expected {
				t.Errorf("categorizeTemperature(%d) = %s; want %s",
					tt.temp, result, tt.expected)
			}
		})
	}
}

func TestFindTodaysPeriod(t *testing.T) {
	service := NewService()

	tests := []struct {
		name     string
		periods  []Period
		expected string
		wantErr  bool
	}{
		{
			name:    "Empty periods",
			periods: []Period{},
			wantErr: true,
		},
		{
			name: "Single daytime period",
			periods: []Period{
				{Name: "Today", IsDaytime: true, ShortForecast: "Sunny"},
			},
			expected: "Sunny",
			wantErr:  false,
		},
		{
			name: "Tonight first, then today",
			periods: []Period{
				{Name: "Tonight", IsDaytime: false, ShortForecast: "Clear"},
				{Name: "Today", IsDaytime: true, ShortForecast: "Sunny"},
			},
			expected: "Sunny",
			wantErr:  false,
		},
		{
			name: "Tonight only",
			periods: []Period{
				{Name: "Tonight", IsDaytime: false, ShortForecast: "Clear"},
			},
			expected: "Clear",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			period, err := service.findTodaysPeriod(tt.periods)
			if (err != nil) != tt.wantErr {
				t.Errorf("findTodaysPeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && period.ShortForecast != tt.expected {
				t.Errorf("findTodaysPeriod() forecast = %v, want %v", period.ShortForecast, tt.expected)
			}
		})
	}
}
