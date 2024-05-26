package main

import (
	"context"
	"os"

	"micro-names/controllers"
	repositories "micro-names/repositories/dynamodb"
	"micro-names/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
)

var repo *repositories.DynamoDBNamesRepository
var service *services.NamesService
var controller *controllers.NamesHttpController
var router = chi.NewRouter()

func main() {
	tableName := os.Getenv("TABLE_NAME")
	sess := session.Must(session.NewSession())
	ddb := dynamodb.New(sess)

	repo = repositories.NewDynamoDBNamesRepository(ddb, tableName)
	service = services.NewNamesService(repo)
	controller = controllers.NewNamesHttpController(service)

	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	router.Get("/dynamo", controller.CreateName)

	return chiadapter.New(router).ProxyWithContext(ctx, req)
}
