package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/TomSED/weather-api/pkg/openweathermap"
	"github.com/TomSED/weather-api/pkg/weatherstack"
	"github.com/aws/aws-lambda-go/events"
)

type WeatherService struct {
	weatherStackClient   WeatherStackClient
	openWeatherMapClient OpenWeatherMapClient
}

func NewWeatherService(weatherStackClient WeatherStackClient, openWeatherMapClient OpenWeatherMapClient) *WeatherService {
	return &WeatherService{
		weatherStackClient:   weatherStackClient,
		openWeatherMapClient: openWeatherMapClient,
	}
}

type GetWeatherResponse struct {
	WindSpeed   int `json:"wind_speed"`
	Temperature int `json:"temperature_degrees"`
}

func (ws *WeatherService) GetWeather(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// For debugging
	requestBytes, _ := json.MarshalIndent(e, "", "    ")
	fmt.Println(string(requestBytes))

	city, exist := e.QueryStringParameters["city"]
	if !exist || city == "" {
		fmt.Println("Missing city in query parameter")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing city in query parameter",
		}, nil
	}

	var weather *GetWeatherResponse
	weatherStackResp, err := ws.weatherStackClient.GetWeather(city)
	if err != nil {
		fmt.Printf("ws.weatherStackClient.GetWeather error: %v\n", err)

		openWeatherMapResp, err := ws.openWeatherMapClient.GetWeather(city)
		if err != nil {
			fmt.Printf("ws.openWeatherMapClient.GetWeather error: %v\n", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Something has gone wrong",
			}, nil
		}
		weather = mapOpenWeatherMapResponse(openWeatherMapResp)

	} else {
		weather = mapWeatherStackResponse(weatherStackResp)
	}

	byt, err := json.Marshal(weather)
	if err != nil {
		fmt.Printf("json.Marshal error: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Something has gone wrong",
		}, nil
	}
	apiResponseBody := string(byt)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       apiResponseBody,
	}, nil
}

func mapWeatherStackResponse(resp *weatherstack.APIResponse) *GetWeatherResponse {
	return &GetWeatherResponse{
		WindSpeed:   resp.Current.WindSpeed,
		Temperature: resp.Current.Temperature,
	}
}

func mapOpenWeatherMapResponse(resp *openweathermap.APIResponse) *GetWeatherResponse {

	temp := int(math.Round(resp.Main.Temp))
	windSpeed := int(math.Round(resp.Wind.Speed))

	return &GetWeatherResponse{
		WindSpeed:   temp,
		Temperature: windSpeed,
	}
}
