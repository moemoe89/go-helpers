package cloudstorage

//go:generate rm -f ./cloudstorage_mock.go
//go:generate mockgen -destination cloudstorage_mock.go -package cloudstorage -mock_names Client=GoMockClient -source cloudstorage.go

import (
	"context"
	"io"
	"time"
)

// CloudFile is a data structure for file in the cloud.
type CloudFile struct {
	// URl is the public URL for the cloud file.
	URL string
}

// Client is an interface for Cloud Storage.
type Client interface {
	// Upload uploads the file to the Cloud Storage given by object
	// and return the CloudFile data structure.
	Upload(ctx context.Context, file io.Reader, object string, expires time.Time) (*CloudFile, error)
	// Delete deletes the given object from Cloud Storage.
	Delete(ctx context.Context, object string) error
}
