package controllers

import (
	"encoding/json"
	"micro-names/services"
	"net/http"
)

type NamesEventsController struct {
	service *services.NamesService
}

func NewNamesEventsController(service *services.NamesService) *NamesEventsController {
	return &NamesEventsController{
		service: service,
	}
}

func (c *NamesEventsController) CreateName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Call controller
	id, err := c.service.CreateName(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	idJSON, err := json.Marshal(map[string]string{"id": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(idJSON)
}
