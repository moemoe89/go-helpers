## tinify

tinify is a Go client library for the Tinify API.
The library provides an implementation of the Image Compressor interface for Tinify service.
The library can be used to compress images using Tinify API.

### Usage

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/moemoe89/go-helpers/pkg/imagecompressor/tinify"
)

func main() {
	// Create a tinify client using an API key.
	client, err := tinify.New(
		// Your API key.
		tinify.WithAPIKey("<API_KEY>"),
		// Optionally customize the http.Client.
		// If not specify, will use http.DefaultClient.
		tinify.WithHTTPClient(
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
	// filename is optional for tinify,
	// leave it as empty string should be fine.
	compressedFile, err := client.Upload(context.Background(), file, "")
	if err != nil {
		panic(err)
	}

	// Print the URL of the compressed file.
	fmt.Println(compressedFile.URL)
}
```

### Options

The New function accepts the following options:

* WithAPIKey: sets the API key for Tinify API.
* WithHTTPClient: sets the HTTPClient interface to make HTTP requests.
