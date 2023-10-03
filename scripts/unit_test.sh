#!/bin/sh
set -o errexit
# Auth
cd auth && go test -race ./... && cd ..
# service/aws/config
cd service/aws/config && go test -race ./... && cd ../../..
# service/aws/sqs
cd service/aws/sqs && go test -race $(go list ./... | grep -v mocks) && cd ../../..
