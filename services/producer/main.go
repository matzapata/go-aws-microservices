package main

import (
	"context"
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

	_, err := ebClient.PutEvents(&eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				EventBusName: aws.String(eventBusName),
				Source:       aws.String("custom.create_name"),
				DetailType:   aws.String("create_name"),
				Detail:       aws.String(fmt.Sprintf(`{"name": "%s"}`, name)),
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
