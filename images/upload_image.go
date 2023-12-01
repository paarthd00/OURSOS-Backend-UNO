package image

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"oursos.com/packages/util"
)

func UploadImage(c echo.Context) error {
	err := godotenv.Load()
	util.CheckError(err)

	credential, err := azblob.NewSharedKeyCredential(os.Getenv("AZURE_ACCOUNT"), os.Getenv("AZURE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", os.Getenv("AZURE_ACCOUNT")))

	serviceURL := azblob.NewServiceURL(*u, p)

	ctx := context.Background() // This example uses a never-expiring context.

	containerURL := serviceURL.NewContainerURL(os.Getenv("AZURE_CONTAINER")) // Container names require lowercase
	json_map := make(map[string]string)

	errEnc := json.NewDecoder(c.Request().Body).Decode(&json_map)
	util.CheckError(errEnc)

	imagefile := json_map["imagefile"]

	// Get the image data and the image format
	dataURLParts := strings.Split(imagefile, ",")
	if len(dataURLParts) != 2 {
		return fmt.Errorf("invalid data URL")
	}

	// Decode the base64 image data
	decodedData, err := base64.StdEncoding.DecodeString(dataURLParts[1])
	if err != nil {
		return err
	}

	// Generate a random hex string for the blob name
	blobName := make([]byte, 16)
	rand.Read(blobName)
	// Append the format as the file extension
	format := strings.Split(strings.Split(dataURLParts[0], "/")[1], ";")[0]
	blobURL := containerURL.NewBlockBlobURL(hex.EncodeToString(blobName) + "." + format)

	// Upload the decoded image data
	_, err = azblob.UploadStreamToBlockBlob(ctx, bytes.NewReader(decodedData), blobURL, azblob.UploadStreamToBlockBlobOptions{
		BufferSize: 4 * 1024 * 1024,
		MaxBuffers: 16,
	})
	if err != nil {
		return err
	}

	// Set the content type of the blob
	_, err = blobURL.SetHTTPHeaders(ctx, azblob.BlobHTTPHeaders{ContentType: "image/" + format}, azblob.BlobAccessConditions{})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"url": blobURL.String()})
}
