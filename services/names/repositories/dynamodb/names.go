package repositories

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type DynamoDBNamesRepository struct {
	DDB       *dynamodb.DynamoDB
	TableName string
}

func NewDynamoDBNamesRepository(ddb *dynamodb.DynamoDB, tableName string) *DynamoDBNamesRepository {
	return &DynamoDBNamesRepository{
		DDB:       ddb,
		TableName: tableName,
	}
}

func (repo *DynamoDBNamesRepository) CreateName(name string) (string, error) {
	id := uuid.New().String()

	input := &dynamodb.PutItemInput{
		TableName: aws.String(repo.TableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"name": {
				S: aws.String(name),
			},
		},
	}

	_, err := repo.DDB.PutItem(input)
	if err != nil {
		return "", err
	}

	return id, nil
}
