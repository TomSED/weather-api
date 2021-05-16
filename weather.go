package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type WeatherService struct {
	weatherStackClient   WeatherStackClient
	openWeatherMapClient OpenWeatherMapClient
}

func NewWeatherService() *WeatherService {
	return &WeatherService{}
}

func (ws *WeatherService) GetWeather(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// For debugging
	requestBytes, _ := json.MarshalIndent(e, "", "    ")
	fmt.Println(string(requestBytes))

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Success",
	}, nil
}
