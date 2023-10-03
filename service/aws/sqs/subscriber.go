package sqs

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	_maxNumberOfMessages = 1
	_waitTimeSeconds     = 1
	_visibilityTimeout   = 1
)

// ErrLengthMessagesInvalid is an error indicating that the number of messages is outside
// the valid range of one to ten.
var ErrLengthMessagesInvalid = errors.New("messages should be no more than ten or less than one")

// Subscriber represents a client for a subscription-based messaging system.
type Subscriber struct {
	sub Consumer
	url string
}

// NewSubscriber creates a new Subscriber instance.
// It returns a pointer to the created Subscriber and an error, if any.
func NewSubscriber(sub Consumer, url string) (*Subscriber, error) {
	if len(url) < 54 {
		return nil, ErrInvalidURL
	}

	return &Subscriber{
		sub: sub,
		url: url,
	}, nil
}

// ReceiveMessage retrieves one or more messages (up to 10).
//
// Using the WaitTimeSeconds parameter enables long-poll support. For more information, see
// Amazon SQS Long Polling (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-long-polling.html)
func (s Subscriber) ReceiveMessage(ctx context.Context, optFns ...func(options *Options)) (*sqs.ReceiveMessageOutput, error) {
	var opt Options
	opt.client = &sqs.ReceiveMessageInput{}
	for _, fn := range optFns {
		fn(&opt)
	}

	if opt.client.MaxNumberOfMessages == 0 {
		opt.client.MaxNumberOfMessages = _maxNumberOfMessages
	}
	if opt.client.WaitTimeSeconds == 0 {
		opt.client.WaitTimeSeconds = _waitTimeSeconds
	}
	if opt.client.VisibilityTimeout == 0 {
		opt.client.VisibilityTimeout = _visibilityTimeout
	}

	output, err := s.sub.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: opt.client.MaxNumberOfMessages,
		WaitTimeSeconds:     opt.client.WaitTimeSeconds,
		VisibilityTimeout:   opt.client.VisibilityTimeout,
	})
	if err != nil {
		return nil, err
	}

	return output, err
}

// DeleteMessages deletes multiple messages from the Amazon SQS queue associated
// with the Subscriber.
//
// The result of the action on each message is reported individually in the
// response (sqs.DeleteMessageBatchOutput). Because the batch request can result
// in a combination of successful and unsuccessful actions, you should check for
// batch errors even when the call returns an HTTP status code of 200 .
func (s Subscriber) DeleteMessages(ctx context.Context, messages ...string) (*sqs.DeleteMessageBatchOutput, error) {
	err := ensureLengthMessages(messages...)
	if err != nil {
		return nil, err
	}

	var entries []types.DeleteMessageBatchRequestEntry
	for i, message := range messages {
		entries = append(entries, types.DeleteMessageBatchRequestEntry{
			Id:            aws.String(fmt.Sprintf("msg%d", i)),
			ReceiptHandle: aws.String(message),
		})
	}

	input := &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(s.url),
		Entries:  entries,
	}

	output, err := s.sub.DeleteMessageBatch(ctx, input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func ensureLengthMessages(messages ...string) error {
	if len(messages) < 1 || len(messages) > 10 {
		return ErrLengthMessagesInvalid
	}

	return nil
}

// Options holds the options for receiving messages from SQS.
type Options struct {
	client *sqs.ReceiveMessageInput
}

// WithMaxNumberOfMessages allows you to configure the MaxNumberOfMessages for use to receive messages.
func WithMaxNumberOfMessages(maxNumberOfMessages int32) func(options *Options) {
	return func(opt *Options) {
		opt.client.MaxNumberOfMessages = maxNumberOfMessages
	}
}

// WithVisibilityTimeout allows you to configure the duration (in seconds)
// that the received messages are hidden from subsequent retrieve requests.
func WithVisibilityTimeout(visibilityTimeout int32) func(options *Options) {
	return func(opt *Options) {
		opt.client.VisibilityTimeout = visibilityTimeout
	}
}

// WithWaitTimeSeconds allows you to configure the WithVisibilityTimeout for use to enables long-poll support.
func WithWaitTimeSeconds(waitTimeSeconds int32) func(options *Options) {
	return func(opt *Options) {
		opt.client.WaitTimeSeconds = waitTimeSeconds
	}
}
