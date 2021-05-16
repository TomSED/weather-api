package openweathermap

const (
	defaultBaseURL = "https://api.openweathermap.org"
)

type Client struct {
	baseURL string
	apiKey  string
}

func NewClient(baseURL string, apiKey string) *Client {

	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}
