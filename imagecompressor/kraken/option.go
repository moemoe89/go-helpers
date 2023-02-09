package kraken

import "net/http"

// Option configures Kraken.
type Option func(k *krakenClient) error

// defaultOptions is a default configuration for kraken.
var defaultOptions = []Option{
	WithHTTPClient(http.DefaultClient),
}

// WithAPIKey returns an option that set the api key.
func WithAPIKey(str string) Option {
	return func(k *krakenClient) error {
		if len(str) == 0 {
			return errFailedSetAPIKey
		}

		k.apiKey = str

		return nil
	}
}

// WithAPISecret returns an option that set the api secret.
func WithAPISecret(str string) Option {
	return func(k *krakenClient) error {
		if len(str) == 0 {
			return errFailedSetAPISecret
		}

		k.apiSecret = str

		return nil
	}
}

// WithHTTPClient returns an option that set the http client.
func WithHTTPClient(httpClient HTTPClient) Option {
	return func(k *krakenClient) error {
		if httpClient == nil {
			return errFailedSetHTTPClient
		}

		k.httpClient = httpClient

		return nil
	}
}
