package kraken

import (
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

			k := &krakenClient{
				apiKey: tt.fields.apiKey,
			}

			err := WithAPIKey(tt.args.value)(k)

			assert.Equal(t, tt.want, k.apiKey)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestWithAPISecret(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		apiSecret string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set api secret value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "api-secret",
				},
				want:    "api-secret",
				wantErr: nil,
			}
		},
		"Failed set api secret value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					apiSecret: "",
				},
				want:    "",
				wantErr: errFailedSetAPISecret,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			k := &krakenClient{
				apiSecret: tt.fields.apiSecret,
			}

			err := WithAPISecret(tt.args.value)(k)

			assert.Equal(t, tt.want, k.apiSecret)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
