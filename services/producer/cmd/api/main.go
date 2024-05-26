package main

import (
	"context"
	"os"
	"producer/controllers"
	"producer/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
)

var controller *controllers.ProducerController
var service *services.ProducerService
var router = chi.NewRouter()
var snsClient *sns.SNS
var topicArn string

func init() {
	topicArn = os.Getenv("TOPIC_ARN")
	sess := session.Must(session.NewSession())
	snsClient = sns.New(sess)

	service = services.NewProducerService(snsClient, topicArn)
	controller = controllers.NewProducerController(service)
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	router.Get("/producer", controller.GetProducer)

	return chiadapter.New(router).ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(HandleRequest)
}
