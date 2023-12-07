# aws config

The service/aws/config module is intended to establish a configuration
with the aws client that will later be used to interact with
the different amazon web services (aws).

Throughout the guide we will see how we can use the go-sdk api
and set configurations for aws-sdk v1, aws-sdk v2 and finally
how to make specific configurations to connect to
[localstack](https://www.localstack.cloud/).

### Install

    go get -u github.com/pomelo-la/go-sdk/service/aws/config

### Configuring aws/config for v2 of aws-sdk-go

```go
package main

import (
	"log"

	"github.com/pomelo-la/go-toolkit/service/aws/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	// Example use to configure sqs client
	// client := sqs.NewFromConfig(*cfg)
	
	// Add your logic
}
```

!!! note
    By default, the client is configured for the **us-east-1** region.

### Configuring aws region

```go
package main

import (
	"log"

	"github.com/pomelo-la/go-toolkit/service/aws/config"
)

func main() {
	sess, err := config.New(config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	// Example use to configure sqs client
	// client := sqs.NewFromConfig(*cfg)

	// Add your logic
}
```

### Usage for localstack

[Localstack](https://www.localstack.cloud/) by default configures
http://localhost:4566 to communicate with local aws services.

Next we will configure the client to be able to access localstack.

```go
package main

import (
	"log"
	
	"github.com/pomelo-la/go-toolkit/service/aws/config"
)

func main() {
	sess, err := config.New(config.WithEndpoint("http://localhost:4566"))
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	// Example use to configure sqs client
	// client := sqs.NewFromConfig(*cfg)

	// Add your logic
}
```

### Configuring client

One design we maintain in our clients is to handle configurations
with _Golang Functional Options Pattern_, this allows flexibility when
configuring clients, here is an example overwriting the
region and client endpoint of the aws/config client.

```go
package main

import (
	"log"

	"github.com/pomelo-la/go-toolkit/service/aws/config"
)

func main() {
	sess, err := config.New(
		config.WithEndpoint("http://mycustomendpoint:1234"),
		config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	// Example use to configure sqs client
	// client := sqs.NewFromConfig(*cfg)

	// Add your logic
}
```

### Configuring aws/config for v1 of aws-sdk-go

```go
package main

import (
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
	
	// Add your logic
}
```

The rest of the API and its behaviour 
is the same for the v1 and v2 client.
