package main

import (
	weatherapi "github.com/TomSED/weather-api"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	ws := weatherapi.NewWeatherService()

	lambda.Start(ws.GetWeather)
}
