//go:generate mockgen -destination ./mocks/uploader.go -package mock -mock_names Uploader=Uploader . Uploader
//go:generate mockgen -destination ./mocks/downloader.go -package mock -mock_names Downloader=Downloader . Downloader
//go:generate mockgen -destination ./mocks/download_upload.go -package mock -mock_names DownloadUploader=DownloadUploader . DownloadUploader

package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// DownloadUploader is an interface that wraps the basic PutObject and GetObject methods.
type DownloadUploader interface {
	Downloader
	Uploader
}

// Uploader is an interface that wraps the basic PutObject methods.
type Uploader interface {
	// PutObject adds an object to an aws s3 bucket.
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// Downloader is an interface that wraps the basic GetObject methods.
type Downloader interface {
	// GetObject retrieves objects from aws S3 bucket.
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}
