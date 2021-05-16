package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/TomSED/weather-api/pkg/openweathermap"
	"github.com/TomSED/weather-api/pkg/weatherstack"
	"github.com/aws/aws-lambda-go/events"
)

// WeatherService provides the lambda handlers for the weather api
type WeatherService struct {
	weatherStackClient   WeatherStackClient
	openWeatherMapClient OpenWeatherMapClient
}

// NewWeatherService creates a new WeatherService
func NewWeatherService(weatherStackClient WeatherStackClient, openWeatherMapClient OpenWeatherMapClient) *WeatherService {
	return &WeatherService{
		weatherStackClient:   weatherStackClient,
		openWeatherMapClient: openWeatherMapClient,
	}
}

// GetWeatherResponse is the struct for the GetWeather api response
type GetWeatherResponse struct {
	WindSpeed   int `json:"wind_speed"`
	Temperature int `json:"temperature_degrees"`
}

// GetWeather is the endpoint for retrieving current temperature and windspeed of a city (via query params city=sydney)
// Weather sources uses weather stack with a failover of open weathermap
func (ws *WeatherService) GetWeather(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// For debugging
	requestBytes, _ := json.MarshalIndent(e, "", "    ")
	fmt.Println(string(requestBytes))

	city, exist := e.QueryStringParameters["city"]
	if !exist || city == "" {
		fmt.Printf(`Missing e.QueryStringParameters["city"]: %v\n`, city)
		return badRequest("Missing city in query parameter"), nil
	}

	var weather *GetWeatherResponse
	weatherStackResp, err := ws.weatherStackClient.GetWeather(city)
	if err != nil {
		fmt.Printf("ws.weatherStackClient.GetWeather error: %v\n", err)

		openWeatherMapResp, err := ws.openWeatherMapClient.GetWeather(city)
		if err != nil {
			fmt.Printf("ws.openWeatherMapClient.GetWeather error: %v\n", err)
			return internalServerError(), nil
		}
		weather = mapOpenWeatherMapResponse(openWeatherMapResp)

	} else {
		weather = mapWeatherStackResponse(weatherStackResp)
	}

	byt, err := json.Marshal(weather)
	if err != nil {
		fmt.Printf("json.Marshal error: %v\n", err)
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
