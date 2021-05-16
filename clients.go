package weatherapi

import (
	"github.com/TomSED/weather-api/pkg/openweathermap"
	"github.com/TomSED/weather-api/pkg/postgres"
	"github.com/TomSED/weather-api/pkg/weatherstack"
)

//go:generate moq -pkg mocks -out mocks/mock_weatherstack_client.go . WeatherStackClient

// WeatherStackClient is an interface for the weather stack api client
type WeatherStackClient interface {
	GetWeather(city string) (*weatherstack.APIResponse, error)
}

//go:generate moq -pkg mocks -out mocks/mock_openweathermap_client.go . OpenWeatherMapClient

// OpenWeatherMapClient is an interface for the open weather map api client
type OpenWeatherMapClient interface {
	GetWeather(city string) (*openweathermap.APIResponse, error)
}

//go:generate moq -pkg mocks -out mocks/mock_postgres_client.go . PostgresClient

// PostgresClient is an interface for the weather database client
type PostgresClient interface {
	InsertWeatherData(*postgres.WeatherData) error
	GetLatestWeatherData() (*postgres.WeatherData, error)
}
