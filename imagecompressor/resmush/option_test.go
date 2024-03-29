package resmush

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

			tp := &reSmushClient{
				httpClient: tt.fields.httpClient,
			}

			err := WithHTTPClient(tt.args.value)(tp)

			assert.Equal(t, tt.want, tp.httpClient)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
