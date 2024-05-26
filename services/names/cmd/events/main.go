package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/matzapata/go-aws-microservices/services/names/controllers"
	repositories "github.com/matzapata/go-aws-microservices/services/names/repositories/dynamodb"
	"github.com/matzapata/go-aws-microservices/services/names/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var repo *repositories.DynamoDBNamesRepository
var service *services.NamesService
var eventsController *controllers.NamesEventsController

type SqsEvent struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}

func init() {
	tableName := os.Getenv("TABLE_NAME")
	sess := session.Must(session.NewSession())
	ddb := dynamodb.New(sess)

	repo = repositories.NewDynamoDBNamesRepository(ddb, tableName)
	service = services.NewNamesService(repo)
	eventsController = controllers.NewNamesEventsController(service)
}

func HandleRequest(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		fmt.Printf("Processing SQS message %s, body: %s\n", record.MessageId, record.Body)

		// Your SQS message processing logic here. Route events to controllers
		var event SqsEvent
		log.Printf("Received SQS message: %s", record.Body)
		err := json.Unmarshal([]byte(record.Body), &event)
		if err != nil {
			return fmt.Errorf("could not unmarshal SQS message body: %v", err)
		}

		// TODO: properly route to controller based on TopicArn
		return eventsController.CreateName(event.Message)
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
