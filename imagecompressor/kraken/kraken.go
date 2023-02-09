package kraken

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/moemoe89/go-helpers/imagecompressor"
)

// apiURL is API URL for Kraken
const apiURL = "https://api.kraken.io/v1/upload"

var (
	// errFailedSetAPIKey is an error message when failed to set api key.
	errFailedSetAPIKey = errors.New("failed to set kraken.api_key")
	// errFailedSetAPISecret is an error message when failed to set api secret.
	errFailedSetAPISecret = errors.New("failed to set kraken.api_secret")
	// errFailedSetHTTPClient is an error message when failed to set http client.
	errFailedSetHTTPClient = errors.New("failed to set kraken.http_client")
	// errInternal is an error message for internal error.
	errInternal = errors.New("internal error")
	// errExternal is an error message for external error.
	errExternal = errors.New("external error")
)

type Data struct {
	Message        string `json:"message"`
	FileName       string `json:"file_name"`
	OriginalSize   int    `json:"original_size"`
	KrakedSize     int    `json:"kraked_size"`
	SavedBytes     int    `json:"saved_bytes"`
	KrakedURL      string `json:"kraked_url"`
	OriginalWidth  int    `json:"original_width"`
	OriginalHeight int    `json:"original_height"`
	KrakedWidth    int    `json:"kraked_width"`
	KrakedHeight   int    `json:"kraked_height"`
	Success        bool   `json:"success"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// krakenClient is a struct for kraken client.
type krakenClient struct {
	apiKey     string
	apiSecret  string
	httpClient HTTPClient
}

func wrapErr(err1 error, err2 error) error {
	return fmt.Errorf("%v: %w", err1, err2)
}

// New returns Image Compressor interface implementations.
func New(opts ...Option) (imagecompressor.Client, error) {
	k := new(krakenClient)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(k); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", wrapErr(err, errInternal))
		}
	}

	return k, nil
}

// Upload uploads the image to the Image Compressor service given by object
// and return the CompressedFile data structure.
// NOTE:
// filename in this service is optional.
func (k *krakenClient) Upload(
	ctx context.Context, file io.Reader, filename string,
) (*imagecompressor.CompressedFile, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	err := writer.WriteField("data", fmt.Sprintf(`{
		"auth":{"api_key": "%s", "api_secret": "%s"}, "wait":true}
	`, k.apiKey, k.apiSecret))
	if err != nil {
		return nil, err
	}

	part, err := writer.CreateFormFile("upload", uuid.New().String())
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

	resp, err := k.httpClient.Do(req)
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
		return nil, fmt.Errorf("failed unmarshal kraken data: %w", wrapErr(err, errInternal))
	}

	if !data.Success {
		return nil, fmt.Errorf("failed to compress image from kraken: %s %w", data.Message, errExternal)
	}

	return &imagecompressor.CompressedFile{
		URL: data.KrakedURL,
	}, nil
}
