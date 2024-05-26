package controllers

import (
	"encoding/json"
	"net/http"
)

type HelloController struct {
}

func NewHelloController() *HelloController {
	return &HelloController{}
}

func (hc *HelloController) GetHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idJSON, err := json.Marshal(map[string]string{"data": "hello!"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(idJSON)
}
