package main

import (
	"context"
	"fmt"
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
var httpController *controllers.NamesHttpController
var router = chi.NewRouter()

func init() {
	tableName := os.Getenv("TABLE_NAME")
	sess := session.Must(session.NewSession())
	ddb := dynamodb.New(sess)

	repo = repositories.NewDynamoDBNamesRepository(ddb, tableName)
	service = services.NewNamesService(repo)
	httpController = controllers.NewNamesHttpController(service)
}

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing API Gateway request: %s\n", event.Body)

	router.Get("/names", httpController.CreateName)

	return chiadapter.New(router).ProxyWithContext(ctx, event)
}

func main() {
	lambda.Start(HandleRequest)
}
