package gcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithBucket(t *testing.T) {
	type args struct {
		value string
	}

	type fields struct {
		bucket string
	}

	type test struct {
		args    args
		fields  fields
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully set bucket value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "bucket",
				},
				want:    "bucket",
				wantErr: nil,
			}
		},
		"Failed set bucket value": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					value: "",
				},
				fields: fields{
					bucket: "example-test",
				},
				want:    "example-test",
				wantErr: errFailedSetBucket,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			g := &gcsClient{
				bucket: tt.fields.bucket,
			}

			err := WithBucket(tt.args.value)(g)

			assert.Equal(t, tt.want, g.bucket)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
