package diskstorage

import (
	"bytes"
	"net/http"
)

// Option configures diskStorage.
type Option func(d *diskStorage) error

// defaultOptions is a default configuration for diskStorage.
var defaultOptions = []Option{
	WithBuffer(&bytes.Buffer{}),
	WithHTTPClient(http.DefaultClient),
}

// WithBuffer returns an option that set the buffer value.
func WithBuffer(b *bytes.Buffer) Option {
	return func(d *diskStorage) error {
		if b == nil {
			return errFailedSetBuffer
		}

		d.buffer = b

		return nil
	}
}

// WithHTTPClient returns an option that set the http client.
func WithHTTPClient(httpClient HTTPClient) Option {
	return func(d *diskStorage) error {
		if httpClient == nil {
			return errFailedSetHTTPClient
		}

		d.httpClient = httpClient

		return nil
	}
}
