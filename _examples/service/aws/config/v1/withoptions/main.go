package main

import (
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/pomelo-la/go-toolkit/service/aws/config"
)

func main() {
	cfg, err := config.New(
		config.WithEndpoint("http://mycustomendpoint:1234"),
		config.WithProfile("my-profile"),
		config.WithRegion("us-west-2"))
	if err != nil {
		return
	}

	s3.New(sess)
}
