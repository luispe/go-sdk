//go:generate mockgen -destination ./mocks/deleter.go -package mock -mock_names Deleter=Deleter . Deleter
//go:generate mockgen -destination ./mocks/downloader.go -package mock -mock_names Downloader=Downloader . Downloader
//go:generate mockgen -destination ./mocks/uploader.go -package mock -mock_names Uploader=Uploader . Uploader
//go:generate mockgen -destination ./mocks/delete_download_upload.go -package mock -mock_names DeleteDownloadUploader=DeleteDownloadUploader . DeleteDownloadUploader

package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// DeleteDownloadUploader is an interface that wraps the basic DeleteObjects, PutObject and GetObject methods.
type DeleteDownloadUploader interface {
	Deleter
	Downloader
	Uploader
}

// Deleter is an interface that wraps the basic DeleteObjects method.
type Deleter interface {
	// DeleteObjects delete object from aws S3 bucket.
	DeleteObjects(ctx context.Context, params *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error)
}

// Downloader is an interface that wraps the basic GetObject method.
type Downloader interface {
	// GetObject retrieves object from aws S3 bucket.
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// Uploader is an interface that wraps the basic PutObject method.
type Uploader interface {
	// PutObject adds an object to an aws s3 bucket.
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}
