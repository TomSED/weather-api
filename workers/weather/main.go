package main

import (
	"os"

	weatherapi "github.com/TomSED/weather-api"
	"github.com/TomSED/weather-api/pkg/openweathermap"
	"github.com/TomSED/weather-api/pkg/postgres"
	"github.com/TomSED/weather-api/pkg/weatherstack"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

func main() {

	logger := logrus.New()
	logger.Out = os.Stdout

	weatherStackClient := weatherstack.NewClient("", os.Getenv("WEATHERSTACK_API_KEY"))
	openWeatherMapClient := openweathermap.NewClient("", os.Getenv("OPENWEATHERMAP_API_KEY"))

	postgresClient, err := postgres.NewClient(os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_USERNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DB_NAME"))
	if err != nil {
		logger.Errorf("postgres.NewClient error: %v", err)
		os.Exit(1)
	}

	ws := weatherapi.NewWeatherService(weatherStackClient, openWeatherMapClient, postgresClient)
	ws.SetLogger(logger)

	lambda.Start(ws.GetWeather)
}
