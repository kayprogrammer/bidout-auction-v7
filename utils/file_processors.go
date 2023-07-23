package utils

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cloudName = os.Getenv("CLOUDINARY_CLOUD_NAME")
var apiKey = os.Getenv("CLOUDINARY_API_KEY")
var apiSecret = os.Getenv("CLOUDINARY_API_SECRET")
var cld, _ = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
var baseFolder = "bidout-auction-v7/"

type SignatureFormat struct {
	PublicId  string `json:"public_id"`
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
}

var ImageExtensions = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/gif":  "gif",
	"image/bmp":  "bmp",
	"image/webp": "webp",
	"image/tiff": "tiff",
}

func retrieveSignatureFromUrl(signedUrl string) string {
	// Parse the signed URL
	parsedURL, err := url.Parse(signedUrl)
	if err != nil {
		log.Fatal("Error parsing URL:", err)
	}

	// Extract query parameters
	queryParams, err := url.ParseQuery(parsedURL.RawQuery)
	if err != nil {
		log.Fatal("Error parsing query parameters:", err)
	}
	return queryParams.Get("signature")
}

func GenerateFileSignature(key string, folder string) SignatureFormat {
	key = fmt.Sprintf("%s%s/%s", baseFolder, folder, key)
	timestamp := time.Now().Unix()
	params := map[string]interface{}{
		"public_id": key,
		"timestamp": timestamp,
	}

	// Convert the params to url.Values
	values := url.Values{}
	for k, v := range params {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	resp, err := api.SignParameters(values, apiSecret)
	log.Fatal("Error signing params: ", err)

	signature := retrieveSignatureFromUrl(resp)
	signatureResp := SignatureFormat{PublicId: key, Signature: signature, Timestamp: timestamp}
	return signatureResp
}

func GenerateFileUrl(key string, folder string, contentType string) string {
	key = fmt.Sprintf("%s%s/%s.%s", baseFolder, folder, key, ImageExtensions[contentType])

	// Generate the Cloudinary URL for the existing resource
	urls, err := cld.Media(key)
	if err != nil {
		log.Println("Error generating Cloudinary URL:", err)
	}
	url, err := urls.String()
	if err != nil {
		log.Println("Error converting to string:", err)
	}
	return url
}

func BoolAddr(b bool) *bool {
	boolVar := b
	return &boolVar
}

func UploadImage(file, key string, folder string) {
	key = fmt.Sprintf("%s%s/%s", baseFolder, folder, key)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cld.Upload.Upload(ctx, file, uploader.UploadParams{PublicID: key, Overwrite: BoolAddr(true), Faces: BoolAddr(true)})
}

// func uploadToCloudinary(filePath string) (string, error) {
// 	// Open the file for upload
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer file.Close()

// 	// Upload the file to Cloudinary
// 	uploadResult, err := uploader.Upload(file, uploader.UploadParams{})
// 	if err != nil {
// 		return "", err
// 	}

// 	// Return the public URL of the uploaded file
// 	return uploadResult.SecureURL, nil
// }

// func main() {
// 	filePath := "/path/to/your/local/file.jpg" // Replace this with your file's actual path
// 	uploadedURL, err := uploadToCloudinary(filePath)
// 	if err != nil {
// 		fmt.Println("Error uploading file:", err)
// 		return
// 	}

// 	fmt.Println("File uploaded successfully. Public URL:", uploadedURL)
// }
