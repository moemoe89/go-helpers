package kraken

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/moemoe89/go-helpers/imagecompressor"

	"github.com/stretchr/testify/assert"
)

type mockedHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

type mockedReadCloser struct {
	err error
}

func (m *mockedReadCloser) Read(p []byte) (n int, err error) {
	return 0, m.err
}

func (m *mockedReadCloser) Close() error {
	return m.err
}

type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}

	type test struct {
		args    args
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully init New": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					opts: []Option{
						WithAPIKey("api-key"),
						WithAPISecret("api-secret"),
					},
				},
				wantErr: nil,
			}
		},
		"Failed init New": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					opts: []Option{
						WithAPIKey(""),
					},
				},
				wantErr: errInternal,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			_, err := New(tt.args.opts...)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpload(t *testing.T) {
	type args struct {
		ctx  context.Context
		file io.Reader
	}

	type fields struct {
		mock *mockedHTTPClient
	}

	type test struct {
		args    args
		fields  fields
		want    *imagecompressor.CompressedFile
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully upload image": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:  context.Background(),
					file: bytes.NewBufferString("test image file content"),
				},
				fields: fields{
					mock: &mockedHTTPClient{
						response: &http.Response{
							Body: io.NopCloser(strings.NewReader(`{
								"file_name": "00be1888-702e-4d86-bfc9-58663d8653da",
								"original_size": 188158,
								"kraked_size": 108425,
								"saved_bytes": 79733,
								"kraked_url": "https://dl.kraken.io/api/f9/d2/7b/test",
								"original_width": 446,
								"original_height": 532,
								"kraked_width": 446,
								"kraked_height": 532,
								"success": true
							}`)),
						},
					},
				},
				want:    &imagecompressor.CompressedFile{URL: "https://dl.kraken.io/api/f9/d2/7b/test"},
				wantErr: nil,
			}
		},
		"Failed to do io.Copy file": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx: context.Background(),
					file: &errorReader{
						err: errInternal,
					},
				},
				want:    nil,
				wantErr: errInternal,
			}
		},
		"Failed to create HTTP request": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:  nil,
					file: bytes.NewBufferString("test image file content"),
				},
				want:    nil,
				wantErr: errInternal,
			}
		},
		"Failed do HTTP request": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:  context.Background(),
					file: bytes.NewBufferString("test image file content"),
				},
				fields: fields{
					mock: &mockedHTTPClient{
						err: errInternal,
					},
				},
				want:    nil,
				wantErr: errInternal,
			}
		},
		"Failed to read body": func(t *testing.T) test {
			t.Helper()

			mockedResp := &http.Response{
				Body: &mockedReadCloser{
					err: errInternal,
				},
			}

			return test{
				args: args{
					ctx:  context.Background(),
					file: bytes.NewBufferString("test image file content"),
				},
				fields: fields{
					mock: &mockedHTTPClient{
						response: mockedResp,
						err:      nil,
					},
				},
				want:    nil,
				wantErr: errInternal,
			}
		},
		"Failed to unmarshal response data": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:  context.Background(),
					file: bytes.NewBufferString("test image file content"),
				},
				fields: fields{
					mock: &mockedHTTPClient{
						response: &http.Response{
							Body: io.NopCloser(strings.NewReader(`invalid json`)),
						},
					},
				},
				want:    nil,
				wantErr: errInternal,
			}
		},
		"Failed to compress due to error from the service": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:  context.Background(),
					file: bytes.NewBufferString("test image file content"),
				},
				fields: fields{
					mock: &mockedHTTPClient{
						response: &http.Response{
							Body: io.NopCloser(strings.NewReader(`{
								"success": false,
								"message": "Only support image type"
							}`)),
						},
					},
				},
				want:    nil,
				wantErr: errExternal,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			sut, err := New(WithHTTPClient(tt.fields.mock))
			assert.NoError(t, err)

			got, err := sut.Upload(tt.args.ctx, tt.args.file, "")
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
