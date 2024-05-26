package controllers

import (
	"net/http"

	"github.com/matzapata/go-aws-microservices/services/names/services"

	"shared/helpers"
)

type NamesHttpController struct {
	service *services.NamesService
}

func NewNamesHttpController(service *services.NamesService) *NamesHttpController {
	return &NamesHttpController{
		service: service,
	}
}

func (c *NamesHttpController) CreateName(w http.ResponseWriter, r *http.Request) {
	// extract data from request
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// business logic
	id, err := c.service.CreateName(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write response
	res := map[string]string{"id": id}
	helpers.WriteJSON(w, http.StatusOK, res)
}
