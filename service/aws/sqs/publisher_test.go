package sqs_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	queue "github.com/pomelo-la/go-toolkit/service/aws/sqs"
	mock "github.com/pomelo-la/go-toolkit/service/aws/sqs/mocks"
)

func TestNewPublisher(t *testing.T) {
	type args struct {
		url string
	}
	type expected struct {
		pub *queue.Publisher
		err error
	}

	tests := []struct {
		name string
		args args
		want expected
	}{
		{
			name: "success",
			args: args{url: "https://sqs.us-east-1.amazonaws.com/012345678901/use1-somequeue"},
			want: expected{pub: mockPublisher(), err: nil},
		},
		{
			name: "error url less characters",
			args: args{url: "https://sqs.us-east-1.amazonaws.com"},
			want: expected{pub: nil, err: queue.ErrInvalidURL},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := queue.NewPublisher(&sqs.Client{}, tt.args.url)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.pub, got)
		})
	}
}

func TestPubSubSendJSONMessage(t *testing.T) {
	type args struct {
		data map[string]any
	}
	type mockPubSub struct {
		output *sqs.SendMessageOutput
		err    error
	}
	type expected struct {
		output *sqs.SendMessageOutput
		err    error
	}

	tests := []struct {
		name       string
		args       args
		mockPubSub *mockPubSub
		want       expected
	}{
		{
			name:       "success",
			args:       args{data: mockMsg()},
			mockPubSub: &mockPubSub{output: &sqs.SendMessageOutput{}, err: nil},
			want:       expected{output: &sqs.SendMessageOutput{}, err: nil},
		},
		{
			name:       "error",
			args:       args{data: mockMsg()},
			mockPubSub: &mockPubSub{output: nil, err: assert.AnError},
			want:       expected{output: nil, err: assert.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			client := mock.NewProducer(ctrl)
			publisher, err := queue.NewPublisher(client, "https://sqs.us-east-1.amazonaws.com/012345678901/use1-somequeue")
			assert.NoError(t, err)

			if tt.mockPubSub != nil {
				client.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(tt.mockPubSub.output, tt.mockPubSub.err).
					Times(1)
			}

			got, err := publisher.SendJSONMessage(ctx, tt.args.data)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.output, got)
		})
	}
}

func TestPublishingSendMessage(t *testing.T) {
	type args struct {
		data map[string]any
	}
	type mockPub struct {
		output *sqs.SendMessageOutput
		err    error
	}
	type expected struct {
		output *sqs.SendMessageOutput
		err    error
	}
	tests := []struct {
		name    string
		args    args
		mockPub *mockPub
		want    expected
	}{
		{
			name:    "success",
			args:    args{data: mockMsg()},
			mockPub: &mockPub{output: &sqs.SendMessageOutput{}, err: nil},
			want:    expected{output: &sqs.SendMessageOutput{}, err: nil},
		},
		{
			name:    "error",
			args:    args{data: mockMsg()},
			mockPub: &mockPub{output: nil, err: assert.AnError},
			want:    expected{output: nil, err: assert.AnError},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			client := mock.NewProducer(ctrl)
			pub, err := queue.NewPublisher(client, "https://sqs.us-east-1.amazonaws.com/012345678901/use1-somequeue")
			assert.NoError(t, err)

			if tt.mockPub != nil {
				client.EXPECT().
					SendMessage(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.mockPub.output, tt.mockPub.err).
					Times(1)
			}

			got, err := pub.SendMessage(ctx, tt.args.data)
			assert.Equal(t, tt.want.output, got)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func mockMsg() map[string]any {
	data := map[string]any{
		"key_int":    1,
		"key_string": "two",
		"key_n":      true,
	}

	return data
}

func mockPublisher() *queue.Publisher {
	pub, _ := queue.NewPublisher(&sqs.Client{}, "https://sqs.us-east-1.amazonaws.com/012345678901/use1-somequeue")
	return pub
}
