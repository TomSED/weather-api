package openweathermap

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWeather(t *testing.T) {

	t.Run("Test Parse Response", func(t *testing.T) {
		cannedResponse := `{"coord":{"lon":151.2073,"lat":-33.8679},"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}],"base":"stations","main":{"temp":289.02,"feels_like":287.52,"temp_min":287.59,"temp_max":290.37,"pressure":1020,"humidity":33},"visibility":10000,"wind":{"speed":5.14,"deg":260},"clouds":{"all":0},"dt":1621130173,"sys":{"type":1,"id":9600,"country":"AU","sunrise":1621111254,"sunset":1621148542},"timezone":36000,"id":2147714,"name":"Sydney","cod":200}`
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
		assert.Equal(t, 5.14, resp.Wind.Speed)
		assert.Equal(t, 289.02, resp.Main.Temp)
	})

	t.Run("Check Request", func(t *testing.T) {

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
		assert.Equal(t, "dummykey", testRequest.URL.Query().Get("appid"))
		assert.Equal(t, "Sydney", testRequest.URL.Query().Get("q"))
		assert.Equal(t, "/data/2.5/weather", testRequest.URL.Path)
	})
}
