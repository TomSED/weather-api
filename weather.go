package weatherapi

import (
	"context"
	"encoding/json"
	"math"

	"github.com/TomSED/weather-api/pkg/openweathermap"
	"github.com/TomSED/weather-api/pkg/weatherstack"
	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
)

// WeatherService provides the lambda handlers for the weather api
type WeatherService struct {
	weatherStackClient   WeatherStackClient
	openWeatherMapClient OpenWeatherMapClient
	logger               *logrus.Logger
}

// NewWeatherService creates a new WeatherService
func NewWeatherService(weatherStackClient WeatherStackClient, openWeatherMapClient OpenWeatherMapClient) *WeatherService {
	return &WeatherService{
		weatherStackClient:   weatherStackClient,
		openWeatherMapClient: openWeatherMapClient,
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

	var weather *GetWeatherResponse
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
		weather = mapOpenWeatherMapResponse(openWeatherMapResp)

	} else {
		weather = mapWeatherStackResponse(weatherStackResp)
	}

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

// mapWeatherStackResponse extracts windspeed & temperature from weatherstack.APIResponse
func mapWeatherStackResponse(resp *weatherstack.APIResponse) *GetWeatherResponse {
	return &GetWeatherResponse{
		WindSpeed:   resp.Current.WindSpeed,
		Temperature: resp.Current.Temperature,
	}
}

// mapOpenWeatherMapResponse extracts windspeed & temperature from openweathermap.APIResponse
func mapOpenWeatherMapResponse(resp *openweathermap.APIResponse) *GetWeatherResponse {

	temp := int(math.Round(resp.Main.Temp))
	windSpeed := int(math.Round(resp.Wind.Speed))

	return &GetWeatherResponse{
		WindSpeed:   temp,
		Temperature: windSpeed,
	}
}
