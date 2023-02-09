package resmush

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
					opts: defaultOptions,
				},
				wantErr: nil,
			}
		},
		"Failed init New": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					opts: []Option{
						WithHTTPClient(nil),
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
		ctx      context.Context
		file     io.Reader
		filename string
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
					ctx:      context.Background(),
					file:     bytes.NewBufferString("test image file content"),
					filename: "test.jpg",
				},
				fields: fields{
					mock: &mockedHTTPClient{
						response: &http.Response{
							Body: io.NopCloser(strings.NewReader(`{
								"dest": "https://resmush.it/test.jpg",
								"dest_size": 123,
								"expires": "2090-06-10",
								"generator": "resmush.it",
								"output": "compressed",
								"percent": 75,
								"src": "https://resmush.it/output",
								"src_size": 256
							}`)),
						},
					},
				},
				want:    &imagecompressor.CompressedFile{URL: "https://resmush.it/test.jpg"},
				wantErr: nil,
			}
		},
		"Failed due to filename empty": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:      nil,
					file:     bytes.NewBufferString("test image file content"),
					filename: "",
				},
				want:    nil,
				wantErr: errExternal,
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
					filename: "test.jpg",
				},
				want:    nil,
				wantErr: errInternal,
			}
		},
		"Failed to create HTTP request": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:      nil,
					file:     bytes.NewBufferString("test image file content"),
					filename: "test.jpg",
				},
				want:    nil,
				wantErr: errInternal,
			}
		},
		"Failed do HTTP request": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:      context.Background(),
					file:     bytes.NewBufferString("test image file content"),
					filename: "test.jpg",
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
					ctx:      context.Background(),
					file:     bytes.NewBufferString("test image file content"),
					filename: "test.jpg",
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
					ctx:      context.Background(),
					file:     bytes.NewBufferString("test image file content"),
					filename: "test.jpg",
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
					ctx:      context.Background(),
					file:     bytes.NewBufferString("test image file content"),
					filename: "test.jpg",
				},
				fields: fields{
					mock: &mockedHTTPClient{
						response: &http.Response{
							Body: io.NopCloser(strings.NewReader(`{
								"error": 400,
								"error_long": "Only support image type"
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

			got, err := sut.Upload(tt.args.ctx, tt.args.file, tt.args.filename)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
