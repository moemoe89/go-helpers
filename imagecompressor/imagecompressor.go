package imagecompressor

//go:generate rm -f ./imagecompressor_mock.go
//go:generate mockgen -destination imagecompressor_mock.go -package imagecompressor -mock_names Client=GoMockClient -source imagecompressor.go

import (
	"context"
	"io"
)

// CompressedFile is a data structure for compressed file in the image compression service.
type CompressedFile struct {
	// URL is the public URL for the compressed file.
	URL string
}

// Client is an interface for Image Compressor service.
type Client interface {
	// Upload uploads the image to the Image Compressor service given by object
	// and return the CompressedFile data structure.
	// NOTE:
	// filename could be optional depending on the image compressor service.
	// Currently only reSmush required filename when uploading to their service.
	Upload(ctx context.Context, file io.Reader, filename string) (*CompressedFile, error)
}
