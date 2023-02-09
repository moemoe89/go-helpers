package diskstorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

//go:generate rm -f ./diskstorage_mock.go
//go:generate mockgen -destination diskstorage_mock.go -package diskstorage -mock_names DiskStorage=GoMockDiskStorage -source diskstorage.go

// DiskStorage defines the methods needed to store and manage files on disk.
type DiskStorage interface {
	// GetBuffer returns the buffer value.
	GetBuffer() *bytes.Buffer
	// WriteFile saves the buffer with given path and permission to disk.
	WriteFile(filePath string, perm os.FileMode) error
	// Write saves the given byte slice to the buffer.
	Write(chunk []byte) error
	// Download downloads file by given target URL and save to file path.
	Download(ctx context.Context, targetURL, filePath string) error
	// Delete removes the file at the given file path.
	Delete(filePath string) error
	// ResetBuffer the buffer value.
	ResetBuffer()
}

// compile time interface implementation check.
var _ DiskStorage = (*diskStorage)(nil)

var (
	// errFailedSetBuffer is an error message when failed to set buffer.
	errFailedSetBuffer = errors.New("failed to set diskstorage.buffer")
	// errFailedSetHTTPClient is an error message when failed to set http client.
	errFailedSetHTTPClient = errors.New("failed to set diskstorage.http_client")
	// errInternal is an error message for internal error.
	errInternal = errors.New("internal error")
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type diskStorage struct {
	buffer     *bytes.Buffer
	httpClient HTTPClient
}

func wrapErr(err1 error, err2 error) error {
	return fmt.Errorf("%v: %w", err1, err2)
}

// New returns a new DiskStorage instance with an empty internal buffer.
func New(opts ...Option) (DiskStorage, error) {
	d := new(diskStorage)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(d); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", wrapErr(err, errInternal))
		}
	}

	return d, nil
}

// GetBuffer returns the buffer value.
func (d *diskStorage) GetBuffer() *bytes.Buffer {
	return d.buffer
}

// WriteFile saves the buffer with given path and permission to disk.
func (d *diskStorage) WriteFile(filePath string, perm os.FileMode) error {
	return os.WriteFile(filePath, d.buffer.Bytes(), perm)
}

// Write saves the given byte slice to the buffer.
func (d *diskStorage) Write(chunk []byte) error {
	_, err := d.buffer.Write(chunk)

	return err
}

// Download downloads file by given target URL and save to file path.
func (d *diskStorage) Download(ctx context.Context, targetURL, filePath string) error {
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", wrapErr(err, errInternal))
	}

	defer func() { _ = out.Close() }()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return fmt.Errorf("failed create http request: %w", wrapErr(err, errInternal))
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send http request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed writing to file: %w", wrapErr(err, errInternal))
	}

	return nil
}

// Delete removes the file at the given file path.
func (d *diskStorage) Delete(filePath string) error {
	return os.Remove(filePath)
}

// ResetBuffer the buffer value.
func (d *diskStorage) ResetBuffer() {
	d.buffer.Reset()
}
