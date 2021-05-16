package openweathermap

const (
	defaultBaseURL = "http://api.weatherstack.com"
)

type Client struct {
	baseURL string
	apiKey  string
}

func NewClient(baseURL string, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}
