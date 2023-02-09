## resmush

resmush is a Go client library for the reSmush API.
The library provides an implementation of the Image Compressor interface for reSmush service.
The library can be used to compress images using reSmush API.

### Usage

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/moemoe89/go-helpers/imagecompressor/resmush"
)

func main() {
	// Create a resmush client.
	client, err := resmush.New(
		// Optionally customize the http.Client.
		// If not specify, will use http.DefaultClient.
		resmush.WithHTTPClient(
			&http.Client{
				Transport:     nil,
				CheckRedirect: nil,
				Jar:           nil,
				Timeout:       0,
			},
		),
	)
	if err != nil {
		panic(err)
	}

	// Open an image file for compression.
	file, err := os.Open("test.jpg")
	if err != nil {
		panic(err)
	}

	// Compress the image and get the compressed file.
	// filename is required for reSmush,
	// need to fill the filename with the extension.
	compressedFile, err := client.Upload(context.Background(), file, "test.jpg")
	if err != nil {
		panic(err)
	}

	// Print the URL of the compressed file.
	fmt.Println(compressedFile.URL)
}
```

### Options

The New function accepts the following options:

* WithHTTPClient: sets the HTTPClient interface to make HTTP requests.
