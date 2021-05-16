package main

import (
	"os"

	weatherapi "github.com/TomSED/weather-api"
	"github.com/TomSED/weather-api/pkg/openweathermap"
	"github.com/TomSED/weather-api/pkg/weatherstack"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

func main() {

	weatherStackClient := weatherstack.NewClient("", os.Getenv("WEATHERSTACK_API_KEY"))
	openWeatherMapClient := openweathermap.NewClient("", os.Getenv("OPENWEATHERMAP_API_KEY"))

	ws := weatherapi.NewWeatherService(weatherStackClient, openWeatherMapClient)

	logger := logrus.New()
	ws.SetLogger(logger)

	lambda.Start(ws.GetWeather)
}
