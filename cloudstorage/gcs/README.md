## gcs

This is a Go package for interacting with Google Cloud Storage (GCS).
It implements the cloudstorage.Client interface and provides methods for uploading and deleting files from GCS.
The package also has methods for generating signed URLs for uploaded files, which can be used for limited time access.

### Usage

To use the package, you first need to create a client:

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/moemoe89/go-helpers/cloudstorage/gcs"
)

func main() {
	ctx := context.Background()

	// Create a new DiskStorage instance with an empty internal buffer.
	client, err := gcs.New(
		ctx,
		// Set the bucket name.
		gcs.WithBucket("test-bucket"),
	)
	if err != nil {
		// handle error
	}

	// Open a file for upload.
	file, err := os.Open("test.jpg")
	if err != nil {
		// handle error
	}

	// Set the filename.
	// Make sure to have it unique depends on the use case.
	object := "test.jpg"

	// Set the expires of file.
	// Leave as time.Time{} if there's no need to the expires time.
	expires := time.Now().Add(time.Hour * 24 * 7)

	// Write data to the buffer.
	cloudFile, err := client.Upload(ctx, file, object, expires)
	if err != nil {
		// handle error
	}

	// Print the URL of the uploaded file.
	fmt.Println(cloudFile.URL)

	// Delete a file on disk.
	err = client.Delete(ctx, object)
	if err != nil {
		// handle error
	}
}
```

### Options

You can customize the client by passing options to the gcs.New function.
The available options are:

* WithBucket: set the name of the GCS bucket to use. If not set, a default bucket will be used.
