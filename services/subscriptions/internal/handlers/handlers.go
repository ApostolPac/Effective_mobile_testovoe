package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"subscriptions/internal/models"
)

type Service interface {
	CreateSub(sub models.Subscription) error
	ReadSub(uuid string) (*models.Subscription, error)
	ReadSubs() ([]models.Subscription, error)
	UpdateSub(sub models.Subscription) error
	DeleteSub(uuid string) error
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

func (h *Handlers) CreateSub(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription
	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.s.CreateSub(sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJSON(w, 200, "Подписка успешно зарегистрирована")
}

func (h *Handlers) ReadSub(w http.ResponseWriter, r *http.Request) {
	uuid := strings.TrimPrefix(r.URL.Path, "/subscriptions/")
	if uuid == "" {
		http.NotFound(w, r)
		return
	}
	sub, error := h.s.ReadSub(uuid)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, 200, sub)
}

func (h *Handlers) UpdateSub(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription
	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.s.UpdateSub(sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJSON(w, 200, "Подписка успешно изменена")
}

func (h *Handlers) DeleteSub(w http.ResponseWriter, r *http.Request) {
	uuid := strings.TrimPrefix(r.URL.Path, "/subscriptions/")
	if uuid == "" {
		http.NotFound(w, r)
		return
	}
	error := h.s.DeleteSub(uuid)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, 200, "Задача успешно удалена")
}

func (h *Handlers) ReadSubs(w http.ResponseWriter, r *http.Request) {

}
func (h *Handlers) ShowMethods(w http.ResponseWriter, r *http.Request) {

}