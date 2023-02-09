## kraken

kraken is a Go client library for the Kraken API.
The library provides an implementation of the Image Compressor interface for Kraken service.
The library can be used to compress images using Kraken API.

### Usage

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/moemoe89/go-helpers/imagecompressor/kraken"
)

func main() {
	// Create a kraken client using an API key.
	client, err := kraken.New(
		// Your API key.
		kraken.WithAPIKey("<API_KEY>"),
		// Your API secret.
		kraken.WithAPISecret("<API_SECRET>"),
		// Optionally customize the http.Client.
		// If not specify, will use http.DefaultClient.
		kraken.WithHTTPClient(
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
	// filename is optional for Kraken,
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

* WithAPIKey: sets the API key for kraken API.
* WithAPISecret: sets the API secret for kraken API.
* WithHTTPClient: sets the HTTPClient interface to make HTTP requests.
