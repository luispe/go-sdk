//go:generate mockgen -destination ./mocks/pub.go -package mock -mock_names Producer=Producer . Producer
//go:generate mockgen -destination ./mocks/sub.go -package mock -mock_names Consumer=Consumer . Consumer

package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Producer is an interface that wraps the basic SendMessage methods.
type Producer interface {
	// SendMessage sends a message using the provided context, parameters, and optional functions.
	// It returns the output of the send message or an error if any.
	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

// Consumer is an interface that wraps the basic ReceiveMessage and DeleteMessage methods.
type Consumer interface {
	// ReceiveMessage one or more messages (up to 10), from the specified queue.
	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)

	DeleteMessage(ctx context.Context,
		params *sqs.DeleteMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}
