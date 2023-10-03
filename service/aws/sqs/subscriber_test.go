package sqs_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	queue "github.com/pomelo-la/go-toolkit/service/aws/sqs"
	mock "github.com/pomelo-la/go-toolkit/service/aws/sqs/mocks"
)

func TestNewSubscriber(t *testing.T) {
	type args struct {
		url string
	}
	type expected struct {
		sub *queue.Subscriber
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
			want: expected{sub: mockSubscriber(), err: nil},
		},
		{
			name: "error url less characters",
			args: args{url: "https://sqs.us-east-1.amazonaws.com"},
			want: expected{sub: nil, err: queue.ErrInvalidURL},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := queue.NewSubscriber(&sqs.Client{}, tt.args.url)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.sub, got)
		})
	}
}

func TestSubscriberReceiveMessage(t *testing.T) {
	type mockPubSub struct {
		output *sqs.ReceiveMessageOutput
		err    error
	}
	type expected struct {
		output *sqs.ReceiveMessageOutput
		err    error
	}

	tests := []struct {
		name       string
		mockPubSub *mockPubSub
		want       expected
	}{
		{
			name:       "success",
			mockPubSub: &mockPubSub{output: &sqs.ReceiveMessageOutput{}, err: nil},
			want:       expected{output: &sqs.ReceiveMessageOutput{}, err: nil},
		},
		{
			name:       "error",
			mockPubSub: &mockPubSub{output: nil, err: assert.AnError},
			want:       expected{output: nil, err: assert.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			client := mock.NewConsumer(ctrl)
			sub, err := queue.NewSubscriber(client, "https://sqs.us-east-1.amazonaws.com/012345678901/use1-somequeue")
			assert.NoError(t, err)

			if tt.mockPubSub != nil {
				client.EXPECT().
					ReceiveMessage(gomock.Any(), gomock.Any()).
					Return(tt.mockPubSub.output, tt.mockPubSub.err).
					Times(1)
			}

			got, err := sub.ReceiveMessage(ctx)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.output, got)
		})
	}
}

func TestSubscriberReceiveMessageWithMaxNumberOfMessages(t *testing.T) {
	type mockPubSub struct {
		output *sqs.ReceiveMessageOutput
		err    error
	}
	type expected struct {
		output *sqs.ReceiveMessageOutput
		err    error
	}

	tests := []struct {
		name       string
		mockPubSub *mockPubSub
		want       expected
	}{
		{
			name:       "success",
			mockPubSub: &mockPubSub{output: &sqs.ReceiveMessageOutput{Messages: make([]types.Message, 7)}, err: nil},
			want:       expected{output: &sqs.ReceiveMessageOutput{Messages: make([]types.Message, 7)}, err: nil},
		},
		{
			name:       "error",
			mockPubSub: &mockPubSub{output: nil, err: assert.AnError},
			want:       expected{output: nil, err: assert.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			client := mock.NewConsumer(ctrl)
			sub, err := queue.NewSubscriber(client, "https://sqs.us-east-1.amazonaws.com/012345678901/use1-somequeue")
			assert.NoError(t, err)

			if tt.mockPubSub != nil {
				client.EXPECT().
					ReceiveMessage(gomock.Any(), gomock.Any()).
					Return(tt.mockPubSub.output, tt.mockPubSub.err).
					Times(1)
			}

			got, err := sub.ReceiveMessage(ctx, queue.WithMaxNumberOfMessages(7))
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.output, got)
			if got != nil {
				assert.Equal(t, len(tt.want.output.Messages), len(got.Messages))
			}
		})
	}
}

func TestSubscriberReceiveMessageWithVisibilityTimeout(t *testing.T) {
	type mockPubSub struct {
		output *sqs.ReceiveMessageOutput
		err    error
	}
	type expected struct {
		output *sqs.ReceiveMessageOutput
		err    error
	}

	tests := []struct {
		name       string
		mockPubSub *mockPubSub
		want       expected
	}{
		{
			name:       "success",
			mockPubSub: &mockPubSub{output: &sqs.ReceiveMessageOutput{}, err: nil},
			want:       expected{output: &sqs.ReceiveMessageOutput{}, err: nil},
		},
		{
			name:       "error",
			mockPubSub: &mockPubSub{output: nil, err: assert.AnError},
			want:       expected{output: nil, err: assert.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			client := mock.NewConsumer(ctrl)
			sub, err := queue.NewSubscriber(client, "https://sqs.us-east-1.amazonaws.com/012345678901/use1-somequeue")
			assert.NoError(t, err)

			if tt.mockPubSub != nil {
				client.EXPECT().
					ReceiveMessage(gomock.Any(), gomock.Any()).
					Return(tt.mockPubSub.output, tt.mockPubSub.err).
					Times(1)
			}

			got, err := sub.ReceiveMessage(ctx, queue.WithVisibilityTimeout(60))
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.output, got)
		})
	}
}

func TestSubscriberReceiveMessageWithWaitTimeSeconds(t *testing.T) {
	type mockPubSub struct {
		output *sqs.ReceiveMessageOutput
		err    error
	}
	type expected struct {
		output *sqs.ReceiveMessageOutput
		err    error
	}

	tests := []struct {
		name       string
		mockPubSub *mockPubSub
		want       expected
	}{
		{
			name:       "success",
			mockPubSub: &mockPubSub{output: &sqs.ReceiveMessageOutput{}, err: nil},
			want:       expected{output: &sqs.ReceiveMessageOutput{}, err: nil},
		},
		{
			name:       "error",
			mockPubSub: &mockPubSub{output: nil, err: assert.AnError},
			want:       expected{output: nil, err: assert.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			client := mock.NewConsumer(ctrl)
			sub, err := queue.NewSubscriber(client, "https://sqs.us-east-1.amazonaws.com/012345678901/use1-somequeue")
			assert.NoError(t, err)

			if tt.mockPubSub != nil {
				client.EXPECT().
					ReceiveMessage(gomock.Any(), gomock.Any()).
					Return(tt.mockPubSub.output, tt.mockPubSub.err).
					Times(1)
			}

			got, err := sub.ReceiveMessage(ctx, queue.WithWaitTimeSeconds(40))
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.output, got)
		})
	}
}

func TestSubscriberDeleteMessage(t *testing.T) {
	type args struct {
		msg string
	}
	type mockPubSub struct {
		output *sqs.DeleteMessageOutput
		err    error
	}
	type expected struct {
		output *sqs.DeleteMessageOutput
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
			args:       args{msg: "some-message"},
			mockPubSub: &mockPubSub{output: &sqs.DeleteMessageOutput{}, err: nil},
			want:       expected{output: &sqs.DeleteMessageOutput{}, err: nil},
		},
		{
			name:       "error",
			args:       args{msg: "some-message"},
			mockPubSub: &mockPubSub{output: nil, err: assert.AnError},
			want:       expected{output: nil, err: assert.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			client := mock.NewConsumer(ctrl)
			sub, err := queue.NewSubscriber(client, "https://sqs.us-east-1.amazonaws.com/012345678901/use1-somequeue")
			assert.NoError(t, err)

			if tt.mockPubSub != nil {
				client.EXPECT().
					DeleteMessage(gomock.Any(), gomock.Any()).
					Return(tt.mockPubSub.output, tt.mockPubSub.err).
					Times(1)
			}

			got, err := sub.DeleteMessage(ctx, tt.args.msg)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.output, got)
		})
	}
}

func mockSubscriber() *queue.Subscriber {
	sub, _ := queue.NewSubscriber(&sqs.Client{}, "https://sqs.us-east-1.amazonaws.com/012345678901/use1-somequeue")
	return sub
}
