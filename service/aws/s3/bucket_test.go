package s3_test

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	bucket "github.com/pomelo-la/go-toolkit/service/aws/s3"
	mock "github.com/pomelo-la/go-toolkit/service/aws/s3/mocks"
)

func TestNewBucket(t *testing.T) {
	type args struct {
		name string
	}
	type expected struct {
		bucket *bucket.Bucket
		err    error
	}

	tests := []struct {
		name string
		args args
		want expected
	}{
		{
			name: "success",
			args: args{name: "some-bucket"},
			want: expected{bucket: mockBucket(), err: nil},
		},
		{
			name: "error name less characters",
			args: args{name: ""},
			want: expected{bucket: nil, err: bucket.ErrInvalidBucketName},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bucket.New(&s3.Client{}, tt.args.name)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.bucket, got)
		})
	}
}

func TestS3PutJSONObject(t *testing.T) {
	type args struct {
		data any
	}
	type mockS3Client struct {
		output *s3.PutObjectOutput
		err    error
	}
	type expected struct {
		output *s3.PutObjectOutput
		err    error
	}

	tests := []struct {
		name         string
		args         args
		mockS3Client *mockS3Client
		want         expected
	}{
		{
			name:         "success",
			args:         args{data: mockData()},
			mockS3Client: &mockS3Client{output: &s3.PutObjectOutput{}, err: nil},
			want:         expected{output: &s3.PutObjectOutput{}, err: nil},
		},
		{
			name:         "error put json object",
			args:         args{data: mockData()},
			mockS3Client: &mockS3Client{output: nil, err: assert.AnError},
			want:         expected{output: nil, err: assert.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			client := mock.NewDownloadUploader(ctrl)
			b, err := bucket.New(client, "some")
			assert.NoError(t, err)

			if tt.mockS3Client != nil {
				client.EXPECT().
					PutObject(gomock.Any(), gomock.Any()).
					Return(tt.mockS3Client.output, tt.mockS3Client.err).
					Times(1)
			}

			got, err := b.PutJSONObject(ctx, tt.args.data)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.output, got)
		})
	}
}

func TestS3PutJSONObjectWithKey(t *testing.T) {
	type args struct {
		data any
	}
	type mockS3Client struct {
		output *s3.PutObjectOutput
		err    error
	}
	type expected struct {
		output *s3.PutObjectOutput
		err    error
	}

	tests := []struct {
		name         string
		args         args
		mockS3Client *mockS3Client
		want         expected
	}{
		{
			name:         "success",
			args:         args{data: mockData()},
			mockS3Client: &mockS3Client{output: &s3.PutObjectOutput{}, err: nil},
			want:         expected{output: &s3.PutObjectOutput{}, err: nil},
		},
		{
			name:         "error put json object",
			args:         args{data: mockData()},
			mockS3Client: &mockS3Client{output: nil, err: assert.AnError},
			want:         expected{output: nil, err: assert.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			client := mock.NewDownloadUploader(ctrl)
			b, err := bucket.New(client, "some")
			assert.NoError(t, err)

			if tt.mockS3Client != nil {
				client.EXPECT().
					PutObject(gomock.Any(), gomock.Any()).
					Return(tt.mockS3Client.output, tt.mockS3Client.err).
					Times(1)
			}

			got, err := b.PutJSONObject(ctx, tt.args.data, bucket.WithKey("my-key"))
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.output, got)
		})
	}
}

func TestS3PutObject(t *testing.T) {
	tempEnvFile := "test.txt"
	f, err := os.Create(tempEnvFile)
	if err != nil {
		t.Fatalf("Failed to create temporary test.txt file: %v", err)
	}
	defer os.Remove(tempEnvFile)
	defer f.Close()

	// Write test environment variables to the temporary .env file
	envData := []byte(`
		some-content
	`)
	_, err = f.Write(envData)
	assert.NoError(t, err)

	type args struct {
		filename string
	}
	type mockS3Client struct {
		output *s3.PutObjectOutput
		err    error
	}
	type expected struct {
		output *s3.PutObjectOutput
		err    error
	}

	tests := []struct {
		name         string
		args         args
		mockS3Client *mockS3Client
		want         expected
	}{
		{
			name:         "success",
			args:         args{filename: "test.txt"},
			mockS3Client: &mockS3Client{output: &s3.PutObjectOutput{}, err: nil},
			want:         expected{output: &s3.PutObjectOutput{}, err: nil},
		},
		{
			name: "error file not found",
			args: args{filename: "not-found.txt"},
			want: expected{output: nil, err: bucket.ErrFileNotFound},
		},
		{
			name:         "error put json object",
			args:         args{filename: "test.txt"},
			mockS3Client: &mockS3Client{output: nil, err: assert.AnError},
			want:         expected{output: nil, err: assert.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			client := mock.NewDownloadUploader(ctrl)
			b, err := bucket.New(client, "some")
			assert.NoError(t, err)

			if tt.mockS3Client != nil {
				client.EXPECT().
					PutObject(gomock.Any(), gomock.Any()).
					Return(tt.mockS3Client.output, tt.mockS3Client.err).
					Times(1)
			}

			got, err := b.PutObject(ctx, tt.args.filename)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.output, got)
		})
	}
}

func TestS3PutObjectWithKey(t *testing.T) {
	tempEnvFile := "test.txt"
	f, err := os.Create(tempEnvFile)
	if err != nil {
		t.Fatalf("Failed to create temporary test.txt file: %v", err)
	}
	defer os.Remove(tempEnvFile)
	defer f.Close()

	// Write test environment variables to the temporary .env file
	envData := []byte(`
		some-content
	`)
	_, err = f.Write(envData)
	assert.NoError(t, err)

	type args struct {
		filename string
	}
	type mockS3Client struct {
		output *s3.PutObjectOutput
		err    error
	}
	type expected struct {
		output *s3.PutObjectOutput
		err    error
	}

	tests := []struct {
		name         string
		args         args
		mockS3Client *mockS3Client
		want         expected
	}{
		{
			name:         "success",
			args:         args{filename: "test.txt"},
			mockS3Client: &mockS3Client{output: &s3.PutObjectOutput{}, err: nil},
			want:         expected{output: &s3.PutObjectOutput{}, err: nil},
		},
		{
			name: "error file not found",
			args: args{filename: "not-found.txt"},
			want: expected{output: nil, err: bucket.ErrFileNotFound},
		},
		{
			name:         "error put json object",
			args:         args{filename: "test.txt"},
			mockS3Client: &mockS3Client{output: nil, err: assert.AnError},
			want:         expected{output: nil, err: assert.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			client := mock.NewDownloadUploader(ctrl)
			b, err := bucket.New(client, "some")
			assert.NoError(t, err)

			if tt.mockS3Client != nil {
				client.EXPECT().
					PutObject(gomock.Any(), gomock.Any()).
					Return(tt.mockS3Client.output, tt.mockS3Client.err).
					Times(1)
			}

			got, err := b.PutObject(ctx, tt.args.filename, bucket.WithKey("my-key"))
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.output, got)
		})
	}
}

func mockData() any {
	return struct {
		Some string
	}{Some: "any"}
}

func mockBucket() *bucket.Bucket {
	b, _ := bucket.New(&s3.Client{}, "some-bucket")
	return b
}
