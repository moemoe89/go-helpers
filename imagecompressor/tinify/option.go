package tinify

import "net/http"

// Option configures tinifyClient.
type Option func(t *tinifyClient) error

// defaultOptions is a default configuration for tinify.
var defaultOptions = []Option{
	WithHTTPClient(http.DefaultClient),
}

// WithAPIKey returns an option that set the api key.
func WithAPIKey(str string) Option {
	return func(t *tinifyClient) error {
		if len(str) == 0 {
			return errFailedSetAPIKey
		}

		t.apiKey = str

		return nil
	}
}

// WithHTTPClient returns an option that set the http client.
func WithHTTPClient(httpClient HTTPClient) Option {
	return func(t *tinifyClient) error {
		if httpClient == nil {
			return errFailedSetHTTPClient
		}

		t.httpClient = httpClient

		return nil
	}
}
