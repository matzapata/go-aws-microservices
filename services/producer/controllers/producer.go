package controllers

import (
	"encoding/json"
	"net/http"
	"producer/services"
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
		w.Header().Set("Content-Type", "application/json")

		idJSON, err := json.Marshal(map[string]string{"error": err.Error()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(idJSON)
	} else {
		w.Header().Set("Content-Type", "application/json")

		idJSON, err := json.Marshal(map[string]string{"data": "OK"})
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Write(idJSON)
	}

}
