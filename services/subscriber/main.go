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

type SnsMessage struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, snsEvent events.SNSEvent) {
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		var message SnsMessage
		err := json.Unmarshal([]byte(snsRecord.Message), &message)
		if err != nil {
			fmt.Printf("Error unmarshaling SNS message: %v\n", err)
			continue
		}

		id := uuid.New().String()
		input := &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]*dynamodb.AttributeValue{
				"id": {
					S: aws.String(id),
				},
				"name": {
					S: aws.String(message.Name),
				},
			},
		}

		_, err = ddb.PutItem(input)
		if err != nil {
			fmt.Printf("Failed to put item: %v\n", err)
		} else {
			fmt.Printf("Successfully processed message: ID = %s, Name = %s\n", id, message.Name)
		}
	}
}

func main() {
	lambda.Start(HandleRequest)
}
