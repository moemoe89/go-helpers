package gcs

// Option configures gcs (Google Cloud Storage).
type Option func(g *gcsClient) error

// defaultOptions is a default configuration for gcs.
var defaultOptions = []Option{
	WithBucket("example-test"),
}

// WithBucket returns an option that set the bucket name.
func WithBucket(str string) Option {
	return func(g *gcsClient) error {
		if len(str) == 0 {
			return errFailedSetBucket
		}

		g.bucket = str

		return nil
	}
}
