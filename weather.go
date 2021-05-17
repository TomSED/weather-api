package weatherapi

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/TomSED/weather-api/pkg/openweathermap"
	"github.com/TomSED/weather-api/pkg/postgres"
	"github.com/TomSED/weather-api/pkg/weatherstack"
	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
)

const (
	CACHE_SECONDS = 3
)

// WeatherService provides the lambda handlers for the weather api
type WeatherService struct {
	weatherStackClient   WeatherStackClient
	openWeatherMapClient OpenWeatherMapClient
	postgresClient       PostgresClient
	logger               *logrus.Logger
}

// NewWeatherService creates a new WeatherService
func NewWeatherService(weatherStackClient WeatherStackClient, openWeatherMapClient OpenWeatherMapClient, postgresClient PostgresClient) *WeatherService {
	return &WeatherService{
		weatherStackClient:   weatherStackClient,
		openWeatherMapClient: openWeatherMapClient,
		postgresClient:       postgresClient,
	}
}

func (ws *WeatherService) SetLogger(logger *logrus.Logger) {
	ws.logger = logger
}

// GetWeatherResponse is the struct for the GetWeather api response
type GetWeatherResponse struct {
	WindSpeed   int `json:"wind_speed"`
	Temperature int `json:"temperature_degrees"`
}

// GetWeather is the endpoint for retrieving current temperature and windspeed of a city (via query params city=sydney)
// Weather sources uses weather stack with a failover of open weathermap
func (ws *WeatherService) GetWeather(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Validate input
	city, exist := e.QueryStringParameters["city"]
	if !exist || city == "" {
		if ws.logger != nil {
			ws.logger.Errorf(`Missing e.QueryStringParameters["city"]: %v\n`, city)
		}
		return badRequest("Missing city in query parameter"), nil
	}

	// Try querying DB
	weatherData, err := ws.postgresClient.GetLatestWeatherData(city)
	if err != nil && ws.logger != nil {
		ws.logger.Errorf("ws.postgresClient.GetLatestWeatherData error: %v\n", err)
	}
	// Check if weather data is up to date
	if err != nil || needsToBeUpdated(weatherData) {
		// Try weatherstack
		weatherStackResp, err := ws.weatherStackClient.GetWeather(city)
		if err != nil {
			if ws.logger != nil {
				ws.logger.Errorf("ws.weatherStackClient.GetWeather error: %v\n", err)
			}

			// Try openweathermap
			openWeatherMapResp, err := ws.openWeatherMapClient.GetWeather(city)
			if err != nil {
				if ws.logger != nil {
					ws.logger.Errorf("ws.openWeatherMapClient.GetWeather error: %v\n", err)
				}
				return internalServerError(), nil
			}
			weatherData = mapOpenWeatherMapResponse(city, openWeatherMapResp)

		} else {
			weatherData = mapWeatherStackResponse(city, weatherStackResp)
		}

		// Update db
		err = ws.postgresClient.InsertWeatherData(weatherData)
		if err != nil && ws.logger != nil {
			ws.logger.Errorf("ws.openWeatherMapClient.GetWeather error: %v\n", err)
			// Non-blocking error, do not need to return a http error, just log error
		}
	}

	// Prepare http response
	var weather *GetWeatherResponse
	weather = mapWeatherData(weatherData)

	// Marshal resp
	byt, err := json.Marshal(weather)
	if err != nil {
		if ws.logger != nil {
			ws.logger.Errorf("json.Marshal error: %v\n", err)
		}
		return internalServerError(), nil
	}
	apiResponseBody := string(byt)

	return success(apiResponseBody), nil
}

func needsToBeUpdated(weatherData *postgres.WeatherData) bool {
	if weatherData == nil {
		return true
	}

	now := time.Now().UTC()
	duration := now.Sub(weatherData.UpdatedDate.UTC())

	if duration > CACHE_SECONDS*time.Second {
		return true
	}

	return false
}

// mapWeatherStackResponse extracts windspeed & temperature from weatherstack.APIResponse
func mapWeatherData(data *postgres.WeatherData) *GetWeatherResponse {
	return &GetWeatherResponse{
		WindSpeed:   data.WindSpeed,
		Temperature: data.Temperature,
	}
}

// mapWeatherStackResponse extracts windspeed & temperature from weatherstack.APIResponse
func mapWeatherStackResponse(city string, resp *weatherstack.APIResponse) *postgres.WeatherData {
	return &postgres.WeatherData{
		DataSource:  "weatherstack",
		City:        city,
		WindSpeed:   resp.Current.WindSpeed,
		Temperature: resp.Current.Temperature,
		UpdatedDate: time.Now().UTC(),
	}
}

// mapOpenWeatherMapResponse extracts windspeed & temperature from openweathermap.APIResponse
func mapOpenWeatherMapResponse(city string, resp *openweathermap.APIResponse) *postgres.WeatherData {

	temp := int(math.Round(resp.Main.Temp))
	windSpeed := int(math.Round(resp.Wind.Speed))

	return &postgres.WeatherData{
		DataSource:  "openweathermap",
		City:        city,
		WindSpeed:   temp,
		Temperature: windSpeed,
		UpdatedDate: time.Now(),
	}
}
