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

var (
	ddb       *dynamodb.DynamoDB
	tableName string
)

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	type MessageContent struct {
		Name string `json:"name"`
	}

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

	for _, message := range sqsEvent.Records {
		var event SqsEvent
		log.Printf("Received SQS message: %s", message.Body)
		err := json.Unmarshal([]byte(message.Body), &event)
		if err != nil {
			return fmt.Errorf("could not unmarshal SQS message body: %v", err)
		}

		var messageContent MessageContent
		err = json.Unmarshal([]byte(event.Message), &messageContent)
		if err != nil {
			fmt.Println("Error:", err)
			return fmt.Errorf("could not unmarshal SQS message: %v", err)
		}

		id := uuid.New().String()
		input := &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]*dynamodb.AttributeValue{
				"id": {
					S: aws.String(id),
				},
				"name": {
					S: aws.String(messageContent.Name),
				},
			},
		}

		_, err = ddb.PutItem(input)
		if err != nil {
			fmt.Printf("Failed to put item: %v\n", err)
		} else {
			fmt.Printf("Successfully processed message: ID = %s, Name = %s\n", id, messageContent.Name)
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
