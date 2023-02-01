package diskstorage

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

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
						WithBuffer(nil),
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

func TestGetBuffer(t *testing.T) {
	type test struct {
		want *bytes.Buffer
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully get buffer": func(t *testing.T) test {
			t.Helper()

			return test{
				want: &bytes.Buffer{},
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			sut, err := New()
			assert.NoError(t, err)

			b := sut.GetBuffer()
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestWriteAndWriteFile(t *testing.T) {
	type args struct {
		chunk []byte
	}

	type test struct {
		args      args
		wantErr   error
		afterFunc func(*testing.T, DiskStorage, []byte)
	}

	defaultAfterFunc := func(t *testing.T, d DiskStorage, data []byte) {
		t.Helper()

		filePath := "test.txt"
		perm := os.FileMode(0644)
		err := d.WriteFile(filePath, perm)
		assert.NoError(t, err)

		defer func() { _ = os.Remove(filePath) }()

		f, err := os.Open(filePath)
		assert.NoError(t, err)

		fileData, err := io.ReadAll(f)
		assert.NoError(t, err)

		if string(fileData) != string(data) {
			t.Errorf("File data does not match expected data. Got %s, expected %s", fileData, data)
		}
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully upload image": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					chunk: []byte("test data"),
				},
				wantErr:   nil,
				afterFunc: defaultAfterFunc,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			sut, err := New()
			assert.NoError(t, err)

			if tt.afterFunc != nil {
				defer tt.afterFunc(t, sut, tt.args.chunk)
			}

			err = sut.Write(tt.args.chunk)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDownload(t *testing.T) {
	type args struct {
		ctx       context.Context
		targetURL string
		filePath  string
	}

	type fields struct {
		mock *mockedHTTPClient
	}

	type test struct {
		args      args
		fields    fields
		wantErr   error
		afterFunc func(*testing.T, string)
	}

	defaultAfterFunc := func(t *testing.T, filePath string) {
		t.Helper()

		defer func() { _ = os.Remove(filePath) }()

		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			t.Errorf("File was not downloaded to the expected location")
		}
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully download file": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:       context.Background(),
					targetURL: "https://www.example.com",
					filePath:  "test_download.txt",
				},
				fields: fields{
					mock: &mockedHTTPClient{
						response: &http.Response{
							Body: io.NopCloser(strings.NewReader(`test download content`)),
						},
					},
				},
				wantErr:   nil,
				afterFunc: defaultAfterFunc,
			}
		},
		"Failed to create file": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:       context.Background(),
					targetURL: "https://www.example.com",
					filePath:  "",
				},
				wantErr: errInternal,
			}
		},
		"Failed to create HTTP request": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:       nil,
					targetURL: "https://www.example.com",
					filePath:  "test_download.txt",
				},
				wantErr: errInternal,
			}
		},
		"Failed do HTTP request": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:       context.Background(),
					targetURL: "https://www.example.com",
					filePath:  "test_download.txt",
				},
				fields: fields{
					mock: &mockedHTTPClient{
						err: errInternal,
					},
				},
				wantErr:   errInternal,
				afterFunc: defaultAfterFunc,
			}
		},
		"Failed to do Copy file": func(t *testing.T) test {
			t.Helper()

			mockedResp := &http.Response{
				Body: &mockedReadCloser{
					err: errInternal,
				},
			}

			return test{
				args: args{
					ctx:       context.Background(),
					targetURL: "https://www.example.com",
					filePath:  "test_download.txt",
				},
				fields: fields{
					mock: &mockedHTTPClient{
						response: mockedResp,
						err:      nil,
					},
				},
				wantErr:   errInternal,
				afterFunc: defaultAfterFunc,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			if tt.afterFunc != nil {
				defer tt.afterFunc(t, tt.args.filePath)
			}

			sut, err := New(WithHTTPClient(tt.fields.mock))
			assert.NoError(t, err)

			err = sut.Download(tt.args.ctx, tt.args.targetURL, tt.args.filePath)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		filePath string
	}

	type test struct {
		args       args
		wantErr    error
		beforeFunc func(*testing.T, string)
		afterFunc  func(*testing.T, string)
	}

	defaultBeforeFunc := func(t *testing.T, filePath string) {
		t.Helper()

		f, err := os.Create(filePath)
		assert.NoError(t, err)

		_ = f.Close()
	}

	defaultAfterFunc := func(t *testing.T, filePath string) {
		t.Helper()

		_, err := os.Stat(filePath)
		if !os.IsNotExist(err) {
			t.Errorf("File was not deleted")
		}
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully upload image": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					filePath: "test_delete.txt",
				},
				wantErr:    nil,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			if tt.beforeFunc != nil {
				tt.beforeFunc(t, tt.args.filePath)
			}

			sut, err := New()
			assert.NoError(t, err)

			if tt.afterFunc != nil {
				defer tt.afterFunc(t, tt.args.filePath)
			}

			err = sut.Delete(tt.args.filePath)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
