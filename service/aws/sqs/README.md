# sqs

Package `sqs` is a small wrapper around https://github.com/aws/aws-sdk-go-v2/tree/main/service/sqs

Welcome to the aws/sqs user guide. Here we will guide you through some
practical examples of how to interact with the API.

# Getting started

```shell
go get github.com/pomelo-la/go-toolkit/service/aws/sqs
```

## Using service aws sqs

The sqs api provides two important structures and associated methods depending on the flow 
with which you wish to interact with sqs, `publisher` and `subscriber`.

### Publisher

Contains the api methods and abstractions necessary to be able to send messages to a queue.


### Subscriber

Contains the api methods and abstractions necessary to be able to 
receive messages from a queue and delete them.
