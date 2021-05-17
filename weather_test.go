package weatherapi_test

import (
	"context"
	"errors"
	"testing"
	"time"

	weatherapi "github.com/TomSED/weather-api"
	"github.com/TomSED/weather-api/mocks"
	"github.com/TomSED/weather-api/pkg/openweathermap"
	"github.com/TomSED/weather-api/pkg/postgres"
	"github.com/TomSED/weather-api/pkg/weatherstack"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestGetWeatherWithDB(t *testing.T) {
	t.Run("If DB success and data is up to date, it should return data from DB", func(t *testing.T) {

		mockPostgresClient := &mocks.PostgresClientMock{
			GetLatestWeatherDataFunc: func(city string) (*postgres.WeatherData, error) {
				return &postgres.WeatherData{
					DataSource:  "datasource",
					Temperature: 1,
					WindSpeed:   2,
					UpdatedDate: time.Now()}, nil
			},
			InsertWeatherDataFunc: func(in1 *postgres.WeatherData) error {
				return nil
			},
		}

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

		mockWeatherService := weatherapi.NewWeatherService(mockWeatherStackClient, mockOpenWeatherMapClient, mockPostgresClient)

		resp, err := mockWeatherService.GetWeather(context.Background(), events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"city": "Sydney",
			},
		})
		if !assert.Nil(t, err) {
			t.Fatal(err)
		}
		assert.Equal(t, 200, resp.StatusCode)
		assert.Len(t, mockPostgresClient.GetLatestWeatherDataCalls(), 1)
		assert.Len(t, mockWeatherStackClient.GetWeatherCalls(), 0)
		assert.Len(t, mockOpenWeatherMapClient.GetWeatherCalls(), 0)
	})

	t.Run("If DB success but data is out of date, it should use weatherstack", func(t *testing.T) {

		mockPostgresClient := &mocks.PostgresClientMock{
			GetLatestWeatherDataFunc: func(city string) (*postgres.WeatherData, error) {
				return &postgres.WeatherData{
					DataSource:  "datasource",
					Temperature: 1,
					WindSpeed:   2,
					UpdatedDate: time.Now().Add(-1 * ((weatherapi.CACHE_SECONDS + 1) * time.Second))}, nil
			},
			InsertWeatherDataFunc: func(in1 *postgres.WeatherData) error {
				return nil
			},
		}

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

		mockWeatherService := weatherapi.NewWeatherService(mockWeatherStackClient, mockOpenWeatherMapClient, mockPostgresClient)

		resp, err := mockWeatherService.GetWeather(context.Background(), events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"city": "Sydney",
			},
		})
		if !assert.Nil(t, err) {
			t.Fatal(err)
		}
		assert.Equal(t, 200, resp.StatusCode)
		assert.Len(t, mockPostgresClient.GetLatestWeatherDataCalls(), 1)
		assert.Len(t, mockWeatherStackClient.GetWeatherCalls(), 1)
		assert.Len(t, mockOpenWeatherMapClient.GetWeatherCalls(), 0)

	})

	t.Run("If DB returns no data, it should use weatherstack", func(t *testing.T) {

		mockPostgresClient := &mocks.PostgresClientMock{
			GetLatestWeatherDataFunc: func(city string) (*postgres.WeatherData, error) {
				return nil, nil
			},
			InsertWeatherDataFunc: func(in1 *postgres.WeatherData) error {
				return nil
			},
		}

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

		mockWeatherService := weatherapi.NewWeatherService(mockWeatherStackClient, mockOpenWeatherMapClient, mockPostgresClient)

		resp, err := mockWeatherService.GetWeather(context.Background(), events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"city": "Sydney",
			},
		})
		if !assert.Nil(t, err) {
			t.Fatal(err)
		}
		assert.Equal(t, 200, resp.StatusCode)
		assert.Len(t, mockPostgresClient.GetLatestWeatherDataCalls(), 1)
		assert.Len(t, mockWeatherStackClient.GetWeatherCalls(), 1)
		assert.Len(t, mockOpenWeatherMapClient.GetWeatherCalls(), 0)
	})

}
func TestGetWeatherWithAPI(t *testing.T) {

	t.Run("If weatherstack succeeds, it should use weather stack", func(t *testing.T) {
		mockPostgresClient := &mocks.PostgresClientMock{
			GetLatestWeatherDataFunc: func(city string) (*postgres.WeatherData, error) {
				return nil, errors.New("db error")
			},
			InsertWeatherDataFunc: func(in1 *postgres.WeatherData) error {
				return nil
			},
		}

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

		mockWeatherService := weatherapi.NewWeatherService(mockWeatherStackClient, mockOpenWeatherMapClient, mockPostgresClient)

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
		mockPostgresClient := &mocks.PostgresClientMock{
			GetLatestWeatherDataFunc: func(city string) (*postgres.WeatherData, error) {
				return nil, errors.New("db error")
			},
			InsertWeatherDataFunc: func(in1 *postgres.WeatherData) error {
				return nil
			},
		}

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

		mockWeatherService := weatherapi.NewWeatherService(mockWeatherStackClient, mockOpenWeatherMapClient, mockPostgresClient)

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
		mockPostgresClient := &mocks.PostgresClientMock{
			GetLatestWeatherDataFunc: func(city string) (*postgres.WeatherData, error) {
				return nil, errors.New("db error")
			},
			InsertWeatherDataFunc: func(in1 *postgres.WeatherData) error {
				return nil
			},
		}

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

		mockWeatherService := weatherapi.NewWeatherService(mockWeatherStackClient, mockOpenWeatherMapClient, mockPostgresClient)

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
		mockPostgresClient := &mocks.PostgresClientMock{
			GetLatestWeatherDataFunc: func(city string) (*postgres.WeatherData, error) {
				return nil, errors.New("db error")
			},
			InsertWeatherDataFunc: func(in1 *postgres.WeatherData) error {
				return nil
			},
		}

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

		mockWeatherService := weatherapi.NewWeatherService(mockWeatherStackClient, mockOpenWeatherMapClient, mockPostgresClient)

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
