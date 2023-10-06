package sqs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

const _groupID = "go-toolkit-publisher"

var (
	// ErrJSONMarshal is an error variable that represents a JSON marshal error.
	ErrJSONMarshal = errors.New("JSON marshal data")
	// ErrInvalidURL is an error indicating that the SQS URL provided is invalid.
	ErrInvalidURL = errors.New("sqs url must no be less of 54 characters")
)

// Publisher is a client for a publishing messaging system.
type Publisher struct {
	pub Producer
	url string
}

// NewPublisher creates a new instance of Publisher.
// It returns a pointer to the created Publisher and an error, if any.
func NewPublisher(pub Producer, url string) (*Publisher, error) {
	if len(url) < 54 {
		return nil, ErrInvalidURL
	}

	return &Publisher{
		pub: pub,
		url: url,
	}, nil
}

// SendJSONMessage sends a JSON message to the Publisher system.
//
// It returns the *sqs.SendMessageOutput and an error, if any.
func (p Publisher) SendJSONMessage(ctx context.Context, data map[string]any) (*sqs.SendMessageOutput, error) {
	msg, err := mapToString(data)
	if err != nil {
		return nil, err
	}
	output, err := p.pub.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(msg),
		QueueUrl:    aws.String(p.url),
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

// SendJSONFifoMessage sends a message to the AWS SQS fifo queue.
//
// Under the hood, converts the map[string]any parameter to a json string.
//
// If the message sending is successful, it returns nil. If there is an error
// during the operation, an error is returned.
func (p Publisher) SendJSONFifoMessage(ctx context.Context, data map[string]any, optFns ...func(options *PublisherOptions)) (*sqs.SendMessageOutput, error) {
	msg, err := mapToString(data)
	if err != nil {
		return nil, err
	}

	var opt PublisherOptions
	opt.client = &sqs.SendMessageInput{}
	for _, fn := range optFns {
		fn(&opt)
	}
	if opt.client.MessageGroupId == nil {
		opt.client.MessageGroupId = aws.String(_groupID)
	}

	if opt.client.MessageDeduplicationId == nil {
		opt.client.MessageDeduplicationId = aws.String(uuid.NewString())
	}

	output, err := p.pub.SendMessage(ctx, &sqs.SendMessageInput{
		MessageDeduplicationId: opt.client.MessageDeduplicationId,
		MessageGroupId:         opt.client.MessageGroupId,
		MessageBody:            aws.String(msg),
		QueueUrl:               aws.String(p.url),
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

// SendMessage sends a message to the AWS SQS queue.
//
// Under the hood, converts the data type any parameter to a string.
//
// If the message sending is successful, it returns nil. If there is an error
// during the operation, an error is returned.
func (p Publisher) SendMessage(ctx context.Context, data any) (*sqs.SendMessageOutput, error) {
	output, err := p.pub.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(fmt.Sprintf("%v", data)),
		QueueUrl:    aws.String(p.url),
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

// SendFifoMessage sends a message to the AWS SQS fifo queue.
//
// Under the hood, converts the data type any parameter to a string.
//
// If the message sending is successful, it returns nil. If there is an error
// during the operation, an error is returned.
func (p Publisher) SendFifoMessage(ctx context.Context, data any, optFns ...func(options *PublisherOptions)) (*sqs.SendMessageOutput, error) {
	var opt PublisherOptions
	opt.client = &sqs.SendMessageInput{}
	for _, fn := range optFns {
		fn(&opt)
	}
	if opt.client.MessageGroupId == nil {
		opt.client.MessageGroupId = aws.String(_groupID)
	}

	if opt.client.MessageDeduplicationId == nil {
		opt.client.MessageDeduplicationId = aws.String(uuid.NewString())
	}
	output, err := p.pub.SendMessage(ctx, &sqs.SendMessageInput{
		MessageDeduplicationId: opt.client.MessageDeduplicationId,
		MessageGroupId:         opt.client.MessageGroupId,
		MessageBody:            aws.String(fmt.Sprintf("%v", data)),
		QueueUrl:               aws.String(p.url),
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

// PublisherOptions holds the options for publisher messages to SQS.
type PublisherOptions struct {
	client *sqs.SendMessageInput
}

// WithGroupID allows you to configure the tag (messageGroupID) to be used to publish messages.
// The tag that specifies that a message belongs to a particular message group.
// More information https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_SendMessage.html
func WithGroupID(groupID string) func(options *PublisherOptions) {
	return func(opt *PublisherOptions) {
		opt.client.MessageGroupId = aws.String(groupID)
	}
}

// WithDeduplicationID allows you to configure the DeduplicationID (messageDeduplicationID) to be used to publish messages.
// The DeduplicationID used for deduplication of sent messages.
// More information https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_SendMessage.html
func WithDeduplicationID(deduplicationID string) func(options *PublisherOptions) {
	return func(opt *PublisherOptions) {
		opt.client.MessageDeduplicationId = aws.String(deduplicationID)
	}
}

func mapToString[T any](data map[string]T) (string, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return "", ErrJSONMarshal
	}

	return string(jsonStr), nil
}
