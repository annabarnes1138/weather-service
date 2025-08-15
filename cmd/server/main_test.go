package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherHandlerMissingParams(t *testing.T) {
	server := NewServer()

	tests := []struct {
		name string
		url  string
	}{
		{"No parameters", "/weather"},
		{"Missing longitude", "/weather?lat=40.7128"},
		{"Missing latitude", "/weather?lng=-74.0060"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.weatherHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusBadRequest)
			}
		})
	}
}

func TestWeatherHandlerInvalidParams(t *testing.T) {
	server := NewServer()

	tests := []struct {
		name string
		url  string
	}{
		{"Invalid latitude", "/weather?lat=invalid&lng=-74.0060"},
		{"Invalid longitude", "/weather?lat=40.7128&lng=invalid"},
		{"Latitude out of range", "/weather?lat=91&lng=-74.0060"},
		{"Longitude out of range", "/weather?lat=40.7128&lng=181"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.weatherHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusBadRequest)
			}
		})
	}
}

func TestWeatherHandlerInvalidMethod(t *testing.T) {
	server := NewServer()

	req, err := http.NewRequest("POST", "/weather?lat=40.7128&lng=-74.0060", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.weatherHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
