package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

var (
	// ErrInvalidBucketName indicates that the bucket name must be between 3 (min) and 63 (max) characters.
	ErrInvalidBucketName = errors.New("bucket name must be between 3 (min) and 63 (max) characters long")
	// ErrFileNotFound is an error that indicates that the specified file was not found.
	ErrFileNotFound = errors.New("the file provided is not found")
)

// Bucket is a client to interact with the object storage.
type Bucket struct {
	client DeleteDownloadUploader
	name   string
}

// New creates a new instance of Bucket.
// It returns a pointer to the created Bucket client and an error, if any.
func New(downloadUploader DeleteDownloadUploader, name string) (*Bucket, error) {
	if len(name) <= 3 || len(name) > 63 {
		return nil, ErrInvalidBucketName
	}
	return &Bucket{client: downloadUploader, name: name}, nil
}

// PutJSONObject adds a go struct JSON object to an aws s3 bucket.
//
// data parameter can be any go struct, under the hood it is converted to json for storage.
func (b Bucket) PutJSONObject(ctx context.Context, data any, optFns ...func(*Options)) (*s3.PutObjectOutput, error) {
	var opt Options
	opt.client = &s3.PutObjectInput{}
	for _, fn := range optFns {
		fn(&opt)
	}
	if opt.client.Key == nil {
		opt.client.Key = aws.String(uuid.NewString())
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}

	output, err := b.client.PutObject(ctx, &s3.PutObjectInput{
		Key:         aws.String(fmt.Sprintf("%s.json", *opt.client.Key)),
		Body:        bytes.NewReader(jsonBytes),
		Bucket:      aws.String(b.name),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

// PutObject adds an object file to an aws s3 bucket.
func (b Bucket) PutObject(ctx context.Context, filename string, optFns ...func(*Options)) (*s3.PutObjectOutput, error) {
	var opt Options
	opt.client = &s3.PutObjectInput{}
	for _, fn := range optFns {
		fn(&opt)
	}
	if opt.client.Key == nil {
		opt.client.Key = aws.String(uuid.NewString())
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, ErrFileNotFound
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("go-toolkit service/aws/s3 unable to close the file reader")
		}
	}()

	output, err := b.client.PutObject(ctx, &s3.PutObjectInput{
		Key:    opt.client.Key,
		Body:   file,
		Bucket: aws.String(b.name),
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DownloadObject retrieves object from aws S3 bucket.
func (b Bucket) DownloadObject(ctx context.Context, objectKey, fileName string) error {
	result, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Key:    aws.String(objectKey),
		Bucket: aws.String(b.name),
	})
	if err != nil {
		return err
	}

	defer func() {
		err := result.Body.Close()
		if err != nil {
			fmt.Println("go-toolkit service/aws/s3 unable to close result GetObject")
		}
	}()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("go-toolkit service/aws/s3 unable to close the file")
		}
	}()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	return nil
}

// DeleteObjects delete objects from aws S3 bucket.
//
// The objectKeys contains a list of up to 1000 keys that you want to delete.
//
// For each objectKeys, Amazon S3 performs a delete
// action and returns the result of that delete, success, or failure, in the
// response (s3.DeleteObjectsOutput). Note that if the object specified in the request is not found, Amazon
// S3 returns the result as deleted.
func (b Bucket) DeleteObjects(ctx context.Context, objectKeys ...string) (*s3.DeleteObjectsOutput, error) {
	var objectIDs []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIDs = append(objectIDs, types.ObjectIdentifier{Key: aws.String(key)})
	}

	output, err := b.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Delete: &types.Delete{Objects: objectIDs},
		Bucket: aws.String(b.name),
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

// Options holds the options for interact to Bucket.
type Options struct {
	client *s3.PutObjectInput
}

// WithKey allows you to configure the key of the object to be named in the bucket.
//
// The object name may or may not contain the extension of the file to be saved in the bucket.
func WithKey(objectKey string) func(*Options) {
	return func(opt *Options) {
		opt.client.Key = aws.String(objectKey)
	}
}
