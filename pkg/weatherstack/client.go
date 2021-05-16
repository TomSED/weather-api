package weatherstack

const (
	defaultBaseURL = "http://api.weatherstack.com"
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
