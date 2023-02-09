package resmush

import "net/http"

// Option configures reSmush.
type Option func(r *reSmushClient) error

// defaultOptions is a default configuration for resmush.
var defaultOptions = []Option{
	WithHTTPClient(http.DefaultClient),
}

// WithHTTPClient returns an option that set the http client.
func WithHTTPClient(httpClient HTTPClient) Option {
	return func(r *reSmushClient) error {
		if httpClient == nil {
			return errFailedSetHTTPClient
		}

		r.httpClient = httpClient

		return nil
	}
}
