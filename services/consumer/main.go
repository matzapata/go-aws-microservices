package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type SQSMessageBody struct {
	Name string `json:"name"`
}

var (
	ddb       *dynamodb.DynamoDB
	tableName string
)

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	type EventDetail struct {
		Name string `json:"name"`
	}

	type CustomEvent struct {
		Version    string      `json:"version"`
		ID         string      `json:"id"`
		DetailType string      `json:"detail-type"`
		Source     string      `json:"source"`
		Account    string      `json:"account"`
		Time       string      `json:"time"`
		Region     string      `json:"region"`
		Resources  []string    `json:"resources"`
		Detail     EventDetail `json:"detail"`
	}

	for _, message := range sqsEvent.Records {
		var event CustomEvent
		log.Printf("Received SQS message: %s", message.Body)
		err := json.Unmarshal([]byte(message.Body), &event)
		if err != nil {
			return fmt.Errorf("could not unmarshal SQS message body: %v", err)
		}

		// Process the event (for demonstration, just print the name)
		fmt.Printf("Processing event: %s\n", event.Detail)

		id := uuid.New().String()
		input := &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]*dynamodb.AttributeValue{
				"id": {
					S: aws.String(id),
				},
				"name": {
					S: aws.String(event.Detail.Name),
				},
			},
		}

		_, err = ddb.PutItem(input)
		if err != nil {
			fmt.Printf("Failed to put item: %v\n", err)
		} else {
			fmt.Printf("Successfully processed message: ID = %s, Name = %s\n", id, event.Detail.Name)
		}
	}

	return nil
}

func init() {
	tableName = os.Getenv("TABLE_NAME")
	sess := session.Must(session.NewSession())
	ddb = dynamodb.New(sess)
}

func main() {
	lambda.Start(HandleRequest)
}
