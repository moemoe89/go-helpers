package resmush

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/moemoe89/go-helpers/imagecompressor"
)

// apiURL is API URL for reSmush
const apiURL = "http://api.resmush.it/ws.php"

var (
	// errFailedSetHTTPClient is an error message when failed to set http client.
	errFailedSetHTTPClient = errors.New("failed to set resmush.http_client")
	// errInternal is an error message for internal error.
	errInternal = errors.New("internal error")
	// errExternal is an error message for external error.
	errExternal = errors.New("external error")
)

type Data struct {
	Src       string `json:"src"`
	Dest      string `json:"dest"`
	SrcSize   int    `json:"src_size"`
	DestSize  int    `json:"dest_size"`
	Percent   int    `json:"percent"`
	Output    string `json:"output"`
	Expires   string `json:"expires"`
	Generator string `json:"generator"`
	Error     int    `json:"error"`
	ErrorLong string `json:"error_long"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// reSmushClient is a struct for reSmush client.
type reSmushClient struct {
	httpClient HTTPClient
}

func wrapErr(err1 error, err2 error) error {
	return fmt.Errorf("%v: %w", err1, err2)
}

// New returns Image Compressor interface implementations.
func New(opts ...Option) (imagecompressor.Client, error) {
	r := new(reSmushClient)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(r); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", wrapErr(err, errInternal))
		}
	}

	return r, nil
}

// Upload uploads the image to the Image Compressor service given by object
// and return the CompressedFile data structure.
// NOTE:
// filename in this service is required.
// reSmush has parameter qlty to choose the quality percentage (0-100).
// In this implementation the qlty not configured so it will use the default value: 92
// reference: https://resmush.it/api
func (r *reSmushClient) Upload(
	ctx context.Context, file io.Reader, filename string,
) (*imagecompressor.CompressedFile, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename for reSmush can't be empty: %w", errExternal)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed create http request: %w", wrapErr(err, errInternal))
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data *Data

	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal resmush data: %w", wrapErr(err, errInternal))
	}

	if data.Error >= 300 {
		return nil, fmt.Errorf("failed to compress image from resmush: %s %w", data.ErrorLong, errExternal)
	}

	return &imagecompressor.CompressedFile{
		URL: data.Dest,
	}, nil
}
