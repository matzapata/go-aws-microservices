package controllers

import (
	"net/http"

	"shared/helpers"
)

type HelloController struct {
}

func NewHelloController() *HelloController {
	return &HelloController{}
}

func (hc *HelloController) GetHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := map[string]string{"data": "hello!"}
	helpers.WriteJSON(w, http.StatusOK, data)
}
