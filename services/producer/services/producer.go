package services

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type ProducerService struct {
	snsClient             *sns.SNS
	CREATE_NAME_TOPIC_ARN string
}

func NewProducerService(snsClient *sns.SNS, CREATE_NAME_TOPIC_ARN string) *ProducerService {
	return &ProducerService{
		snsClient:             snsClient,
		CREATE_NAME_TOPIC_ARN: CREATE_NAME_TOPIC_ARN,
	}
}

func (s *ProducerService) PublishCreateNameEvent(name string) error {
	message := map[string]string{"name": name}
	messageJson, _ := json.Marshal(message)

	_, err := s.snsClient.Publish(&sns.PublishInput{
		TopicArn: aws.String(s.CREATE_NAME_TOPIC_ARN),
		Message:  aws.String(string(messageJson)),
	})

	return err
}
