package weatherapi

import (
	"github.com/TomSED/weather-api/pkg/openweathermap"
	"github.com/TomSED/weather-api/pkg/weatherstack"
)

type WeatherStackClient interface {
	GetWeather(city string) (*weatherstack.APIResponse, error)
}

type OpenWeatherMapClient interface {
	GetWeather(city string) (*openweathermap.APIResponse, error)
}
