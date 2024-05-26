package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/matzapata/go-aws-microservices/services/names/services"
)

type NamesEventsController struct {
	service *services.NamesService
}

func NewNamesEventsController(service *services.NamesService) *NamesEventsController {
	return &NamesEventsController{
		service: service,
	}
}

func (c *NamesEventsController) CreateName(eventMessage string) error {
	type MessageContent struct {
		Name string `json:"name"`
	}

	var messageContent MessageContent
	err := json.Unmarshal([]byte(eventMessage), &messageContent)
	if err != nil {
		fmt.Println("Error:", err)
		return fmt.Errorf("could not unmarshal SQS message: %v", err)
	}

	_, err = c.service.CreateName(messageContent.Name)
	if err != nil {
		return err
	}

	return nil
}
