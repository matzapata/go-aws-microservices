
# TODO: 

- SQS and EventBus?
- Testing
- Split code in controllers, services and so on
- Create helpers for writing json and use DTOS
- Staging env

controllers - business logic
router - map endpoints to controllers, read and write responses

# SQS and EventBus

```ts
import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as eventbridge from 'aws-cdk-lib/aws-events';
import * as eventtargets from 'aws-cdk-lib/aws-events-targets';
import * as lambdatargets from 'aws-cdk-lib/aws-lambda-event-sources';

export class CdkEventbridgeSqsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Create the SQS queue
    const queue = new sqs.Queue(this, 'EventQueue');

    // Create producer Lambda
    const producerLambda = new lambda.Function(this, 'ProducerLambda', {
      runtime: lambda.Runtime.NODEJS_14_X,
      handler: 'producer.handler',
      code: lambda.Code.fromAsset('lambda/producer'),
    });

    // Create consumer Lambda
    const consumerLambda = new lambda.Function(this, 'ConsumerLambda', {
      runtime: lambda.Runtime.NODEJS_14_X,
      handler: 'consumer.handler',
      code: lambda.Code.fromAsset('lambda/consumer'),
    });

    // Add the SQS event source to the consumer Lambda
    consumerLambda.addEventSource(new lambdatargets.SqsEventSource(queue, {
        batchSize: 1 // Process one message per Lambda invocation
    }));

    // Create EventBridge rule for 'create_name' event
    const createNameRule = new eventbridge.Rule(this, 'CreateNameRule', {
      eventPattern: {
        source: ['custom.create_name'],
      },
    });

    // Create EventBridge rule for 'update_name' event
    const updateNameRule = new eventbridge.Rule(this, 'UpdateNameRule', {
      eventPattern: {
        source: ['custom.update_name'],
      },
    });

    // Add the SQS queue as the target for both EventBridge rules
    createNameRule.addTarget(new eventtargets.SqsQueue(queue));
    updateNameRule.addTarget(new eventtargets.SqsQueue(queue));

    // Grant permissions for the producer Lambda to send messages to EventBridge
    producerLambda.addToRolePolicy(new cdk.aws_iam.PolicyStatement({
      actions: ['events:PutEvents'],
      resources: ['*'], // Adjust this to be more restrictive in a real application
    }));

    // Grant the producer Lambda permissions to send messages to the SQS queue
    queue.grantSendMessages(producerLambda);
  }
}
```

## Consumer

```golang
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type SQSMessageBody struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		var body SQSMessageBody
		if err := json.Unmarshal([]byte(message.Body), &body); err != nil {
			log.Printf("Error unmarshalling message body: %v", err)
			continue
		}

		// Process the message
		fmt.Printf("Processing name: %s\n", body.Name)

		// Here you can add your business logic
		// For example, store the name in a database, call another service, etc.
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
```

## Producer

```golang
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

var (
	eventBusName string
	ebClient     *eventbridge.EventBridge
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := request.QueryStringParameters["name"]
	if name == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing 'name' query parameter"}, nil
	}

	detail := map[string]string{"name": name}
	detailJson, _ := json.Marshal(detail)
	_, err := ebClient.PutEvents(&eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				EventBusName: aws.String(eventBusName),
				Source:       aws.String("custom.create_name"),
				DetailType:   aws.String("create_name"),
				Detail:       aws.String(string(detailJson)),
			},
		},
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Failed to publish to EventBridge: %v", err)}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Event published successfully"}, nil
}

func main() {
	eventBusName = os.Getenv("EVENT_BUS_NAME")
	sess := session.Must(session.NewSession())
	ebClient = eventbridge.New(sess)

	lambda.Start(HandleRequest)
}
```