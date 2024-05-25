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
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	snsClient *sns.SNS
	topicArn  string
)

func init() {
	topicArn = os.Getenv("TOPIC_ARN")
	sess := session.Must(session.NewSession())
	snsClient = sns.New(sess)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := request.QueryStringParameters["name"]
	if name == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing 'name' query parameter"}, nil
	}

	message := map[string]string{"name": name}
	messageJson, _ := json.Marshal(message)
	_, err := snsClient.Publish(&sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(string(messageJson)),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Failed to publish to SNS: %v", err)}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Event published successfully"}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
