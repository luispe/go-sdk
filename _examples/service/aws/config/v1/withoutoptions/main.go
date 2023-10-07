package main

import (
	"log"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/pomelo-la/go-toolkit/service/aws/config"
)

func main() {
	sess, err := config.NewV1()
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	s3.New(sess)
}
