# Weather Service

A simple HTTP server that provides weather forecasts using the National Weather Service API.

## Project Structure

```
weather-service/
├── cmd/
│   └── server/
│       ├── main.go          # HTTP server and handlers
│       └── main_test.go     # Server tests
├── pkg/
│   └── weather/
│       ├── client.go        # NWS API client
│       ├── service.go       # Business logic
│       └── service_test.go  # Service tests
├── go.mod                   # Go module definition
└── README.md               # This file
```

## Features

- Accepts latitude and longitude coordinates
- Returns today's short forecast (e.g., "Partly Cloudy")
- Categorizes temperature as "hot", "cold", or "moderate"
- Uses the official National Weather Service API as data source
- Follows Go project layout conventions
- Separation of concerns with proper package structure

## Requirements

- Go 1.21 or later

## Building and Running

1. Clone or download the project files
2. Build the project from the root directory:
   ```bash
   go build -o weather-service ./cmd/server
   ```
3. Run the server:
   ```bash
   ./weather-service
   ```
   
   Or run directly:
   ```bash
   go run ./cmd/server
   ```

The server will start on port 8080 by default.

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run tests for specific packages:
```bash
go test ./pkg/weather
go test ./cmd/server
```

## API Usage

### Get Weather Forecast

**Endpoint:** `GET /weather`

**Parameters:**
- `lat` (required): Latitude coordinate (-90 to 90)
- `lng` (required): Longitude coordinate (-180 to 180)

**Example Request:**
```bash
curl "http://localhost:8080/weather?lat=40.7128&lng=-74.0060"
```

**Example Response:**
```json
{
  "forecast": "Partly Cloudy",
  "temperature": "moderate",
  "location": "40.7128, -74.0060"
}
```

## Temperature Categorization

- **Hot:** 80°F and above
- **Cold:** 50°F and below  
- **Moderate:** 51-79°F

## Example Coordinates

- New York City: `lat=40.7128&lng=-74.0060`
- Los Angeles: `lat=34.0522&lng=-118.2437`
- Chicago: `lat=41.8781&lng=-87.6298`
- Miami: `lat=25.7617&lng=-80.1918`

## Architecture

The project follows Go conventions with clear separation of concerns:

### `cmd/server/`
- **main.go**: HTTP server setup, routing, and request handling
- **main_test.go**: HTTP handler tests

### `pkg/weather/`
- **client.go**: NWS API client with HTTP communication logic
- **service.go**: Business logic for weather processing and temperature categorization
- **service_test.go**: Unit tests for business logic

This structure provides:
- **Testability**: Each layer can be tested independently
- **Maintainability**: Clear boundaries between HTTP handling and business logic
- **Reusability**: The `pkg/weather` package could be imported by other applications
- **Go Conventions**: Follows standard Go project layout

## Error Handling

The service handles various error conditions:
- Invalid coordinates
- Missing parameters
- NWS API failures
- Network timeouts

## Design Notes and Shortcuts

This is a simplified implementation focused on the core requirements. In a production environment, you would want to consider:

### Current Shortcuts:
1. **No authentication/rate limiting** - The NWS API is free but has rate limits
2. **Basic error handling** - Could be more granular with specific error types
3. **No caching** - Each request hits the NWS API directly
4. **Simple temperature categorization** - Uses basic thresholds
5. **No configuration** - Port and thresholds are hardcoded
6. **No logging framework** - Uses basic Go log package
7. **No metrics/monitoring** - No observability beyond basic logging
8. **Basic HTTP client** - Uses `http.Get` instead of custom client with timeout/User-Agent

### Production Considerations:
1. **Caching layer** - Redis/Memcached to cache forecast responses
2. **Configuration management** - Environment variables or config files
3. **Structured logging** - JSON logs with correlation IDs
4. **Metrics and monitoring** - Prometheus metrics, health checks
5. **Graceful shutdown** - Handle SIGTERM/SIGINT properly
6. **Rate limiting** - Protect against abuse
7. **Input validation** - More robust coordinate and parameter validation
8. **Circuit breaker** - Handle NWS API outages gracefully
9. **Retry logic** - Exponential backoff for API failures
10. **HTTPS/TLS** - Secure communications
11. **Dependency injection** - For better testability with mock clients
12. **Context usage** - For request tracing and cancellation
13. **Custom HTTP client** - With proper timeouts, User-Agent headers, and connection pooling

## NWS API Integration

The service uses a two-step process with the National Weather Service API:
1. Convert lat/lng to NWS grid coordinates using `/points/{lat},{lng}`
2. Fetch forecast using `/gridpoints/{office}/{gridX},{gridY}/forecast`

This follows the official NWS API workflow and ensures we get the most accurate forecast data for the specific location.

## Testing the Service

Test the service with various coordinates:

```bash
# Test valid coordinates
curl "http://localhost:8080/weather?lat=40.7128&lng=-74.0060"

# Test invalid coordinates
curl "http://localhost:8080/weather?lat=invalid&lng=-74.0060"

# Test missing parameters
curl "http://localhost:8080/weather?lat=40.7128"
```

## Package Documentation

Generate and view package documentation:
```bash
go doc ./pkg/weather
go doc ./pkg/weather.Service
```