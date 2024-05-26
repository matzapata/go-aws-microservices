package controllers

import (
	"encoding/json"
	"micro-names/services"
	"net/http"
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
	w.Header().Set("Content-Type", "application/json")
	idJSON, err := json.Marshal(map[string]string{"id": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(idJSON)
}
