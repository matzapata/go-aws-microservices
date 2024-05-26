package main

import (
	"context"

	"hello/controllers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
)

var controller *controllers.HelloController
var router = chi.NewRouter()

func main() {
	controller = controllers.NewHelloController()

	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	router.Get("/hello", controller.GetHello)

	return chiadapter.New(router).ProxyWithContext(ctx, req)
}
