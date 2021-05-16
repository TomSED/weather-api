package weatherapi_test

import (
	"context"
	"errors"
	"testing"

	weatherapi "github.com/TomSED/weather-api"
	"github.com/TomSED/weather-api/mocks"
	"github.com/TomSED/weather-api/pkg/openweathermap"
	"github.com/TomSED/weather-api/pkg/weatherstack"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestGetWeather(t *testing.T) {

	t.Run("If weatherstack succeeds, it should use weather stack", func(t *testing.T) {
		mockWeatherStackClient := &mocks.WeatherStackClientMock{
			GetWeatherFunc: func(city string) (*weatherstack.APIResponse, error) {
				resp := &weatherstack.APIResponse{}
				resp.Current.Temperature = 10
				resp.Current.WindSpeed = 11
				return resp, nil
			},
		}

		mockOpenWeatherMapClient := &mocks.OpenWeatherMapClientMock{
			GetWeatherFunc: func(city string) (*openweathermap.APIResponse, error) {
				return nil, errors.New("openweathermap error")
			},
		}

		mockWeatherService := weatherapi.NewWeatherService(mockWeatherStackClient, mockOpenWeatherMapClient)

		resp, err := mockWeatherService.GetWeather(context.Background(), events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"city": "Sydney",
			},
		})
		if !assert.Nil(t, err) {
			t.Fatal(err)
		}
		assert.Equal(t, 200, resp.StatusCode)
		assert.Len(t, mockWeatherStackClient.GetWeatherCalls(), 1)
		assert.Len(t, mockOpenWeatherMapClient.GetWeatherCalls(), 0)
	})

	t.Run("If weatherstack fails, it should use openweather map", func(t *testing.T) {
		mockWeatherStackClient := &mocks.WeatherStackClientMock{
			GetWeatherFunc: func(city string) (*weatherstack.APIResponse, error) {
				return nil, errors.New("weatherstack error")

			},
		}

		mockOpenWeatherMapClient := &mocks.OpenWeatherMapClientMock{
			GetWeatherFunc: func(city string) (*openweathermap.APIResponse, error) {
				resp := &openweathermap.APIResponse{}
				resp.Main.Temp = 10
				resp.Wind.Speed = 11
				return resp, nil
			},
		}

		mockWeatherService := weatherapi.NewWeatherService(mockWeatherStackClient, mockOpenWeatherMapClient)

		resp, err := mockWeatherService.GetWeather(context.Background(), events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"city": "Sydney",
			},
		})
		if !assert.Nil(t, err) {
			t.Fatal(err)
		}
		assert.Equal(t, 200, resp.StatusCode)
		assert.Len(t, mockWeatherStackClient.GetWeatherCalls(), 1)
		assert.Len(t, mockOpenWeatherMapClient.GetWeatherCalls(), 1)
	})

	t.Run("If both data sources fail, it should return a 500 response", func(t *testing.T) {
		mockWeatherStackClient := &mocks.WeatherStackClientMock{
			GetWeatherFunc: func(city string) (*weatherstack.APIResponse, error) {
				return nil, errors.New("weatherstack error")
			},
		}

		mockOpenWeatherMapClient := &mocks.OpenWeatherMapClientMock{
			GetWeatherFunc: func(city string) (*openweathermap.APIResponse, error) {
				return nil, errors.New("openweathermap error")
			},
		}

		mockWeatherService := weatherapi.NewWeatherService(mockWeatherStackClient, mockOpenWeatherMapClient)

		resp, err := mockWeatherService.GetWeather(context.Background(), events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"city": "Sydney",
			},
		})
		if !assert.Nil(t, err) {
			t.Fatal(err)
		}
		assert.Equal(t, 500, resp.StatusCode)
		assert.Len(t, mockWeatherStackClient.GetWeatherCalls(), 1)
		assert.Len(t, mockOpenWeatherMapClient.GetWeatherCalls(), 1)
	})

	t.Run("If no city in query provided, it should return a 400 error", func(t *testing.T) {
		mockWeatherStackClient := &mocks.WeatherStackClientMock{
			GetWeatherFunc: func(city string) (*weatherstack.APIResponse, error) {
				resp := &weatherstack.APIResponse{}
				resp.Current.Temperature = 10
				resp.Current.WindSpeed = 11
				return resp, nil
			},
		}

		mockOpenWeatherMapClient := &mocks.OpenWeatherMapClientMock{
			GetWeatherFunc: func(city string) (*openweathermap.APIResponse, error) {
				resp := &openweathermap.APIResponse{}
				resp.Main.Temp = 10
				resp.Wind.Speed = 11
				return resp, nil
			},
		}

		mockWeatherService := weatherapi.NewWeatherService(mockWeatherStackClient, mockOpenWeatherMapClient)

		// city == ""
		resp, err := mockWeatherService.GetWeather(context.Background(), events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"city": "",
			},
		})
		if !assert.Nil(t, err) {
			t.Fatal(err)
		}
		assert.Equal(t, 400, resp.StatusCode)
		assert.Len(t, mockWeatherStackClient.GetWeatherCalls(), 0)
		assert.Len(t, mockOpenWeatherMapClient.GetWeatherCalls(), 0)

		// city not provided
		resp, err = mockWeatherService.GetWeather(context.Background(), events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{},
		})
		if !assert.Nil(t, err) {
			t.Fatal(err)
		}
		assert.Equal(t, 400, resp.StatusCode)
		assert.Len(t, mockWeatherStackClient.GetWeatherCalls(), 0)
		assert.Len(t, mockOpenWeatherMapClient.GetWeatherCalls(), 0)
	})
}
