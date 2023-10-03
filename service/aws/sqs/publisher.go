package sqs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// ErrJSONMarshal is an error variable that represents a JSON marshal error.
var (
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

func mapToString[T any](data map[string]T) (string, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return "", ErrJSONMarshal
	}

	return string(jsonStr), nil
}
