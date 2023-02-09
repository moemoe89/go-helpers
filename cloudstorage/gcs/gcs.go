package gcs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/moemoe89/go-helpers/cloudstorage"

	"cloud.google.com/go/storage"
)

const (
	// publicHost is a public host for Google Cloud Storage e.g.
	// https://storage.googleapis.com/example-test/test.jpg
	publicHost = "https://storage.googleapis.com"
)

var (
	// errFailedSetBucket is an error message when failed to set bucket.
	errFailedSetBucket = errors.New("failed to set gcs.bucket")
	// errInternal is an error message for internal error.
	errInternal = errors.New("internal error")
	// errExternal is an error message for external error.
	errExternal = errors.New("external error")
)

// gcsClient is a struct for gcs client.
type gcsClient struct {
	*storage.Client

	bucket string
}

func wrapErr(err1 error, err2 error) error {
	return fmt.Errorf("%v: %w", err1, err2)
}

// New returns Cloud Storage interface implementations.
func New(ctx context.Context, opts ...Option) (cloudstorage.Client, error) {
	g := new(gcsClient)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(g); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", wrapErr(err, errInternal))
		}
	}

	var err error

	g.Client, err = storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// Upload uploads the file to the Cloud Storage given by object
// and return the public url of the file.
func (g *gcsClient) Upload(
	ctx context.Context, file io.Reader, object string, expires time.Time,
) (*cloudstorage.CloudFile, error) {
	wc := g.Bucket(g.bucket).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return nil, fmt.Errorf("failed to copy file %s to bucket %s: %w", object, g.bucket, err)
	}

	if err := wc.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer %s to bucket %s: %w", object, g.bucket, err)
	}

	cloudFile := &cloudstorage.CloudFile{
		URL: g.buildURL(object),
	}

	// immediately do return if expires time not configured.
	if expires.IsZero() {
		return cloudFile, nil
	}

	url, err := g.signedURL(object, expires)
	if err != nil {
		return nil, err
	}

	cloudFile.URL = url

	return cloudFile, nil
}

// Delete deletes the given object from Cloud Storage.
func (g *gcsClient) Delete(ctx context.Context, object string) error {
	err := g.Bucket(g.bucket).Object(object).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete file %s on bucket %s: %w", object, g.bucket, err)
	}

	return nil
}

// signedURL signed the object from cloud storage with expires time.
func (g *gcsClient) signedURL(object string, expires time.Time) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  http.MethodGet,
		Expires: expires,
	}

	url, err := g.Bucket(g.bucket).SignedURL(object, opts)
	if err != nil {
		return "", fmt.Errorf("failed to signed url for object %s in bucket %s: %w", object, g.bucket, err)
	}

	return url, nil
}

// buildURL builds the object URL from cloud storage.
func (g *gcsClient) buildURL(object string) string {
	return fmt.Sprintf("%s/%s/%s", publicHost, g.bucket, object)
}
