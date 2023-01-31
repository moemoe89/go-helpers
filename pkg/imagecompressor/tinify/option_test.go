package tinify

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithAPIKey(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		apiKey string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set api key value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "api-key",
				},
				want:    "api-key",
				wantErr: nil,
			}
		},
		"Failed set api key value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					apiKey: "",
				},
				want:    "",
				wantErr: errFailedSetAPIKey,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			tp := &tinifyClient{
				apiKey: tt.fields.apiKey,
			}

			err := WithAPIKey(tt.args.value)(tp)

			assert.Equal(t, tt.want, tp.apiKey)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithHTTPClient(t *testing.T) {
	type args struct {
		value HTTPClient
	}

	type fields struct {
		httpClient HTTPClient
	}

	type test struct {
		args    args
		fields  fields
		want    HTTPClient
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set http client value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: http.DefaultClient,
				},
				want:    http.DefaultClient,
				wantErr: nil,
			}
		},
		"Failed set http client value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: nil,
				},
				fields: fields{
					httpClient: http.DefaultClient,
				},
				want:    http.DefaultClient,
				wantErr: errFailedSetHTTPClient,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			tp := &tinifyClient{
				httpClient: tt.fields.httpClient,
			}

			err := WithHTTPClient(tt.args.value)(tp)

			assert.Equal(t, tt.want, tp.httpClient)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
