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
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

var (
	ddb       *dynamodb.DynamoDB
	tableName string
)

func init() {
	tableName = os.Getenv("TABLE_NAME")
	sess := session.Must(session.NewSession())
	ddb = dynamodb.New(sess)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := request.QueryStringParameters["name"]
	if name == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing 'name' query parameter"}, nil
	}

	id := uuid.New().String()

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"name": {
				S: aws.String(name),
			},
		},
	}

	_, err := ddb.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Failed to put item: %v", err)}, err
	}

	response := map[string]string{
		"id":   id,
		"name": name,
	}
	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(body)}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
