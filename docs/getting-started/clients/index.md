# Getting started for Clients

## Before started

One of the main concepts for working with go-sdk is that each
module/package works separately, this makes administration easier
and for the client it is a better experience since it only
downloads the transitive dependencies of a single module.

## 1. Working with service/aws/config

The service/aws/config module is intended to establish a configuration
with the aws client that will later be used to interact with
the different amazon web services (aws).

### Install

    go get -u github.com/pomelo-la/go-sdk/service/aws/config

### Configuring aws/config for v1 of aws-sdk

```go
package main

import (
	"fmt"
	"log"

	"github.com/pomelo-la/go-toolkit/service/aws/config"
)

func main() {
	sess, err := config.NewV1()
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	// Example use to configure s3 client
	// client := s3.New(sess)
}
```

!!! note 
    By default, the client is configured for the **us-east-1** region.

Please see the full documentation [here](../../packages.md)