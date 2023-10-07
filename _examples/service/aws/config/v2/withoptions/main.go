package main

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/pomelo-la/go-toolkit/service/aws/config"
)

func main() {
	cfg, err := config.New(
		config.WithEndpoint("http://mycustomendpoint:1234"),
		config.WithProfile("my-profile"),
		config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	sqs.NewFromConfig(cfg)
}
