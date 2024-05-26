package controllers

import (
	"net/http"

	"github.com/matzapata/go-aws-microservices/services/producer/services"

	"shared/helpers"
)

type ProducerController struct {
	service *services.ProducerService
}

func NewProducerController(service *services.ProducerService) *ProducerController {
	return &ProducerController{service: service}
}

func (c *ProducerController) GetProducer(w http.ResponseWriter, r *http.Request) {
	// extract name from query param
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// publish the event and write response
	err := c.service.PublishCreateNameEvent(name)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
	} else {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"data": "OK"})
	}

}
