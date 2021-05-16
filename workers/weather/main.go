package main

import (
	"os"

	weatherapi "github.com/TomSED/weather-api"
	"github.com/TomSED/weather-api/pkg/openweathermap"
	"github.com/TomSED/weather-api/pkg/weatherstack"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	weatherStackClient := weatherstack.NewClient("", os.Getenv("WEATHERSTACK_API_KEY"))
	openWeatherMapClient := openweathermap.NewClient("", os.Getenv("OPENWEATHERMAP_API_KEY"))

	ws := weatherapi.NewWeatherService(weatherStackClient, openWeatherMapClient)

	lambda.Start(ws.GetWeather)
}
