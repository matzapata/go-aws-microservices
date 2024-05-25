package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	repositories "dynamo/repositories/dynamodb"
	services "dynamo/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var repo *repositories.DynamoDBNamesRepository
var service *services.NamesService

func main() {
	tableName := os.Getenv("TABLE_NAME")
	sess := session.Must(session.NewSession())
	ddb := dynamodb.New(sess)

	repo = repositories.NewDynamoDBNamesRepository(ddb, tableName)
	service = services.NewNamesService(repo)

	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := request.QueryStringParameters["name"]
	if name == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing 'name' query parameter"}, nil
	}

	id, err := service.CreateName(name)
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
