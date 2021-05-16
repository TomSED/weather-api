package weatherstack

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWeather(t *testing.T) {

	t.Run("Check Request", func(t *testing.T) {
		cannedResponse := `{"request":{"type":"City","query":"Sydney, Australia","language":"en","unit":"m"},"location":{"name":"Sydney","country":"Australia","region":"New South Wales","lat":"-33.883","lon":"151.217","timezone_id":"Australia/Sydney","localtime":"2021-05-16 11:37","localtime_epoch":1621165020,"utc_offset":"10.0"},"current":{"observation_time":"01:37 AM","temperature":15,"weather_code":113,"weather_icons":["https://assets.weatherstack.com/images/wsymbols01_png_64/wsymbol_0001_sunny.png"],"weather_descriptions":["Sunny"],"wind_speed":24,"wind_degree":250,"wind_dir":"WSW","pressure":1021,"precip":0,"humidity":39,"cloudcover":0,"feelslike":14,"uv_index":4,"visibility":10,"is_day":"yes"}}`
		testAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Canned Response
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cannedResponse))
		}))

		client := &Client{
			baseURL: testAPI.URL,
			apiKey:  "dummykey",
		}

		resp, err := client.GetWeather("Sydney")
		if !assert.Nil(t, err) {
			t.Fatal(err)
		}
		assert.Equal(t, 24, resp.Current.WindSpeed)
		assert.Equal(t, 15, resp.Current.Temperature)
	})

	t.Run("Test Parse Response", func(t *testing.T) {

		var testRequest *http.Request
		testAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Canned Response
			w.WriteHeader(200)
			testRequest = req
		}))

		client := &Client{
			baseURL: testAPI.URL,
			apiKey:  "dummykey",
		}

		_, _ = client.GetWeather("Sydney")
		if !assert.NotNil(t, testRequest) {
			t.Fatal()
		}
		assert.Equal(t, http.MethodGet, testRequest.Method)
		assert.Equal(t, "dummykey", testRequest.URL.Query().Get("access_key"))
		assert.Equal(t, "Sydney", testRequest.URL.Query().Get("query"))
	})
}
