package gcs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts []Option
	}

	type test struct {
		args    args
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Failed init New": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx: context.Background(),
					opts: []Option{
						WithBucket(""),
					},
				},
				wantErr: errInternal,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			_, err := New(tt.args.ctx, tt.args.opts...)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
