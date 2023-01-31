package tinify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/moemoe89/go-helpers/pkg/imagecompressor"
)

// apiURL is API URL for Tinify
const apiURL = "https://api.tinify.com/shrink"

var (
	// errFailedSetAPIKey is an error message when failed to set api key.
	errFailedSetAPIKey = errors.New("failed to set tinify.api_key")
	// errFailedSetHTTPClient is an error message when failed to set http client.
	errFailedSetHTTPClient = errors.New("failed to set tinify.http_client")
	// errInternal is an error message for internal error.
	errInternal = errors.New("internal error")
	// errExternal is an error message for external error.
	errExternal = errors.New("external error")
)

type Data struct {
	Input struct {
		Size int    `json:"size"`
		Type string `json:"type"`
	} `json:"input"`
	Output struct {
		Size   int     `json:"size"`
		Type   string  `json:"type"`
		Width  int     `json:"width"`
		Height int     `json:"height"`
		Ratio  float64 `json:"ratio"`
		URL    string  `json:"url"`
	} `json:"output"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// tinifyClient is a struct for Tinify client.
type tinifyClient struct {
	apiKey string

	httpClient HTTPClient
}

func wrapErr(err1 error, err2 error) error {
	return fmt.Errorf("%v: %w", err1, err2)
}

// New returns Image Compressor interface implementations.
func New(opts ...Option) (imagecompressor.Client, error) {
	t := new(tinifyClient)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(t); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", wrapErr(err, errInternal))
		}
	}

	return t, nil
}

// Upload uploads the image to the Image Compressor service given by object
// and return the CompressedFile data structure.
// NOTE:
// filename in this service is optional.
func (t *tinifyClient) Upload(
	ctx context.Context, file io.Reader, filename string,
) (*imagecompressor.CompressedFile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, file)
	if err != nil {
		return nil, fmt.Errorf("failed create http request: %w", wrapErr(err, errInternal))
	}

	req.SetBasicAuth("api", t.apiKey)

	resp, err := t.httpClient.Do(req)
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
		return nil, fmt.Errorf("failed unmarshal tinify data: %w", wrapErr(err, errInternal))
	}

	if data.Error != "" {
		return nil, fmt.Errorf("failed to compress image from tinify: %s. %s %w", data.Error, data.Message, errExternal)
	}

	return &imagecompressor.CompressedFile{
		URL: data.Output.URL,
	}, nil
}
