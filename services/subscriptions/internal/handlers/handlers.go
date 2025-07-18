package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"subscriptions/internal/models"
)

type Service interface {
	CreateSub(w http.ResponseWriter, r *http.Request)(error)
	ReadSub(w http.ResponseWriter, r *http.Request)(*models.Subscription, error)
	ReadSubs(w http.ResponseWriter, r *http.Request)([]models.Subscription, error)
	UpdateSub(w http.ResponseWriter, r *http.Request)(error)
	DeleteSub(w http.ResponseWriter, r *http.Request)(error)
}

type Handlers struct {
	s Service
}

func NewHandler(a Service) *Handlers {
	return &Handlers{
		s: a,
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

func CreateSub() {

}

func ReadSub() {
	
}

func UpdateSub() {

}

func DeleteSub() {

}
