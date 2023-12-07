# aws sqs

The sqs api provides two important structures and associated methods depending on the flow
with which you wish to interact with sqs, `publisher` and `subscriber`.

Publisher contains the api methods and abstractions necessary to be able to send messages to a queue.

Subscriber contains the api methods and abstractions necessary to be able to
receive messages from a queue and delete them.

### Install

    go get -u github.com/pomelo-la/go-sdk/service/aws/sqs

## Publisher

### Working with standard queue

First let's create the client to interact with the queue.

```go
package main

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pomelo-la/go-toolkit/service/aws/config"

	queue "github.com/pomelo-la/go-toolkit/service/aws/sqs"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	fifoPublisher, err := queue.NewPublisher(sqs.NewFromConfig(*cfg), "https://sqs.us-east-1.amazonaws.com/<YOUR_ACCOUNT_ID>/<YOUR_SQS>.fifo")
	if err != nil {
		log.Fatalf("%v", err)
	}
}
```

Next let's send a JSON message

```go hl_lines="25-36"
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pomelo-la/go-toolkit/service/aws/config"

	queue "github.com/pomelo-la/go-toolkit/service/aws/sqs"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	stdPublisher, err := queue.NewPublisher(sqs.NewFromConfig(*cfg), "https://sqs.<AWS_REGION>.amazonaws.com/<YOUR_ACCOUNT_ID>/<YOUR_SQS>")
	if err != nil {
		log.Fatalf("%v", err)
	}
	
	data := map[string]any{
		"key_1": "some value",
		"key_2": false,
		"key_n": 1184,
	}
	
	output, err := stdPublisher.SendJSONMessage(context.Background(), data)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(output)
}
```

Send a simple message

```go hl_lines="25-30"
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pomelo-la/go-toolkit/service/aws/config"

	queue "github.com/pomelo-la/go-toolkit/service/aws/sqs"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	stdPublisher, err := queue.NewPublisher(sqs.NewFromConfig(*cfg), "https://sqs.<AWS_REGION>.amazonaws.com/<YOUR_ACCOUNT_ID>/<YOUR_SQS>")
	if err != nil {
		log.Fatalf("%v", err)
	}
	
	output, err := stdPublisher.SendMessage(context.Background(), "my-message")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(output)
}
```

### Working wit fifo queue

First let's create the client to interact with the queue.

!!! note

    Please find further information about aws fifo queue
    [here](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/welcome.html) 

```go
package main

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pomelo-la/go-toolkit/service/aws/config"

	queue "github.com/pomelo-la/go-toolkit/service/aws/sqs"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	fifoPublisher, err := queue.NewPublisher(sqs.NewFromConfig(*cfg), "https://sqs.<AWS_REGION>.amazonaws.com/<YOUR_ACCOUNT_ID>/<YOUR_SQS>.fifo")
	if err != nil {
		log.Fatalf("%v", err)
	}
}
```

Let's send a JSON message

```go hl_lines="25-36"
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pomelo-la/go-toolkit/service/aws/config"

	queue "github.com/pomelo-la/go-toolkit/service/aws/sqs"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	fifoPublisher, err := queue.NewPublisher(sqs.NewFromConfig(*cfg), "https://sqs.us-east-1.amazonaws.com/<YOUR_ACCOUNT_ID>/<YOUR_SQS>.fifo")
	if err != nil {
		log.Fatalf("%v", err)
	}
	
	data := map[string]any{
		"key_1": "some value",
		"key_2": false,
		"key_n": 1184,
	}
	
	output, err := fifoPublisher.SendJSONFifoMessage(context.Background(), data)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(output)
}
```

Now let's send the same JSON message, but with 
deduplicationID and groupID.

```go hl_lines="33-34"
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pomelo-la/go-toolkit/service/aws/config"

	queue "github.com/pomelo-la/go-toolkit/service/aws/sqs"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	fifoPublisher, err := queue.NewPublisher(sqs.NewFromConfig(*cfg), "https://sqs.us-east-1.amazonaws.com/<YOUR_ACCOUNT_ID>/<YOUR_SQS>.fifo")
	if err != nil {
		log.Fatalf("%v", err)
	}

	data := map[string]any{
		"key_1": "some value",
		"key_2": false,
		"key_n": 1184,
	}
	
	output, err := fifoPublisher.SendJSONFifoMessage(context.Background(),
		data,
		queue.WithGroupID("my-group-id"),
		queue.WithDeduplicationID("abc-123"),
	)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(output)
}
```

!!! tip

    WithGroupID and WithDeduplicationID are a _functional option pattern_ 
    if not sent internally they will be sent with the following defaults 
    GroupID="go-toolkit-publisher" and DeduplicationID="NewUUID"

Finally, the API offers SendFifoMessage which sends a string message 
and has the same features as SendJSONFifoMessage.

---

## Subscriber

### Receive messages

ReceiveMessage retrieves one or more messages (up to 10).
Using the WaitTimeSeconds parameter enables long-poll support. 
For more information, see
[Amazon SQS Long Polling](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-long-polling.html)

First let's create the client to interact with the queue.

```go
package main

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pomelo-la/go-toolkit/service/aws/config"
	queue "github.com/pomelo-la/go-toolkit/service/aws/sqs"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	subscriber, err := queue.NewSubscriber(sqs.NewFromConfig(*cfg), "https://sqs.<AWS_REGION>.amazonaws.com/<YOUR_ACCOUNT_ID>/<YOUR_SQS>.fifo")
	if err != nil {
		log.Fatalf("%v", err)
	}
}
```

Next let's receive messages from a queue

```go hl_lines="24-29"
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pomelo-la/go-toolkit/service/aws/config"
	queue "github.com/pomelo-la/go-toolkit/service/aws/sqs"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("unable to load aws SDK config, %v", err)
	}

	subscriber, err := queue.NewSubscriber(sqs.NewFromConfig(*cfg), "https://sqs.us-east-1.amazonaws.com/<YOUR_ACCOUNT_ID>/<YOUR_SQS>.fifo")
	if err != nil {
		log.Fatalf("%v", err)
	}

	output, err := subscriber.ReceiveMessage(context.Background())
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(output)
}
```

The following methods can be used to configure the way in 
which you wish to consume messages.

    // WithMaxNumberOfMessages allows you to configure the 
    // MaxNumberOfMessages for use to receive messages.
    // default = 1
    WithMaxNumberOfMessages(maxNumberOfMessages int32)

    // WithVisibilityTimeout allows you to configure the duration (in seconds)
    // that the received messages are hidden from subsequent retrieve requests.
    // default = 1
    WithVisibilityTimeout(visibilityTimeout int32)
    
    // WithWaitTimeSeconds allows you to configure the WithWaitTimeSeconds for use
    // to enables long-poll support.
    // default = 1
    WithWaitTimeSeconds(waitTimeSeconds int32)

### Delete message(s)

DeleteMessages deletes multiple messages from the Amazon SQS queue associated
with the Subscriber.

The result of the action on each message is reported individually in the
response (sqs.DeleteMessageBatchOutput). Because the batch request can result
in a combination of successful and unsuccessful actions, you should check for
batch errors even when the call returns an HTTP status code of 200.