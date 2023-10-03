// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pomelo-la/go-toolkit/service/aws/sqs (interfaces: Producer)
//
// Generated by this command:
//
//	mockgen -destination ./mocks/pub.go -package mock -mock_names Producer=Producer . Producer
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	gomock "go.uber.org/mock/gomock"
)

// Producer is a mock of Producer interface.
type Producer struct {
	ctrl     *gomock.Controller
	recorder *ProducerMockRecorder
}

// ProducerMockRecorder is the mock recorder for Producer.
type ProducerMockRecorder struct {
	mock *Producer
}

// NewProducer creates a new mock instance.
func NewProducer(ctrl *gomock.Controller) *Producer {
	mock := &Producer{ctrl: ctrl}
	mock.recorder = &ProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Producer) EXPECT() *ProducerMockRecorder {
	return m.recorder
}

// SendMessage mocks base method.
func (m *Producer) SendMessage(arg0 context.Context, arg1 *sqs.SendMessageInput, arg2 ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendMessage", varargs...)
	ret0, _ := ret[0].(*sqs.SendMessageOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *ProducerMockRecorder) SendMessage(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*Producer)(nil).SendMessage), varargs...)
}
