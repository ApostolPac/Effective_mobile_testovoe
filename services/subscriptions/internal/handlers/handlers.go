package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"subscriptions/internal/models"
	"time"
)

const (
	dataStructError = "неправильный формат данных"
	missedLinkError = "ссылка на подписку отсутствует"
)

type Service interface {
	CreateSub(sub models.Subscription) error                 // Метод для создания записи.
	ReadSub(id int) (*models.Subscription, error)            // Метод для чтения записи по её id.
	ReadSubs(userId string) ([]models.Subscription, error)   // Метод для чтения среза записей для конкретного пользователя.
	UpdateSub(sub models.Subscription) error                 // Метод для обновления записей методом Update.
	DeleteSub(id int) error                                  // Метод для удаления записи о подписке.
	ShowSubscSum(startPeriod time.Time, EndPeriod time.Time) // Метод для получения сум подписок, для начала работы нужно -
	// отправить период внутри которого будем искать записи о подписках
	ShowMethods() error // Метод для показа всех доступных путей и методов.
}

type Handlers struct {
	s Service
}

func NewHandler(a Service) *Handlers {
	return &Handlers{
		s: a,
	}
}

// Функция для сериализации в JSON любого типа данных с помощью пустого интерфейса и его запись в http.ResponseWriter

func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, "ошибка при кодировании данных", http.StatusBadRequest)
		log.Print(err.Error(), " writeJSON method")
		return err
	}

	return nil
}

// Функция для получения id записи

func getSubId(w http.ResponseWriter, r *http.Request) (subIdi int, err error) {
	id := strings.TrimPrefix(r.URL.Path, "/subscriptions/")
	if id == "" {
		http.Error(w, missedLinkError, http.StatusBadRequest)
		log.Print("ошибка при извлечении id - отсутствует id")
		return 0, fmt.Errorf("id подписки отсутствует")
	}

	subId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Ошибка в формате id, должен быть числом", http.StatusInternalServerError)
		log.Print(err.Error(), " getSubId method")
		return 0, fmt.Errorf("id подписки не является числом")
	}

	return subId, nil
}

// Хендлер для создания записи о подписке

func (h *Handlers) CreateSub(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription

	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		http.Error(w, dataStructError, http.StatusBadRequest)
		log.Print(err.Error())
		return
	}

	err = h.s.CreateSub(sub)
	if err != nil {
		http.Error(w, "ошибка при создании записи о подписке", http.StatusInternalServerError)
		log.Print(err.Error(), " CreateSub method")
	}

	err = writeJSON(w, http.StatusOK, "запись о подписке успешно зарегистрирована")
	if err != nil {
		log.Print(err.Error(), " CreateSub method")
		return
	}
}

// Хендлер для чтения записи о подписке

func (h *Handlers) ReadSub(w http.ResponseWriter, r *http.Request) {
	uuid, err := getSubId(w, r)
	if err != nil {
		log.Print(err.Error(), " ReadSub method")
		return
	}

	sub, err := h.s.ReadSub(uuid)
	// if err == sql.ErrNoRows {
	// 	http.Error(w, "запись о подписке не найдена", http.StatusNotFound)
	// 	return
	// }

	if err != nil {
		http.Error(w, "ошибка при получении данных по записе о подписке", http.StatusInternalServerError)
		log.Print(err.Error(), " ReadSub method")
		return
	}

	writeJSON(w, http.StatusOK, sub)
}

// Хендлер для обновления записи о подписке

func (h *Handlers) UpdateSub(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription

	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		http.Error(w, dataStructError, http.StatusBadRequest)
		log.Print(err.Error(), " UpdateSub method")
		return
	}

	err = h.s.UpdateSub(sub)
	if err != nil {
		http.Error(w, "ошибка при обновлении записи о подписке", http.StatusInternalServerError)
		log.Print(err.Error(), " UpdateSub method")
	}

	writeJSON(w, http.StatusOK, "подписка успешно изменена")
}

// Хендлер для удаления записи о подписке

func (h *Handlers) DeleteSub(w http.ResponseWriter, r *http.Request) {
	//удаляем часть пути для получения только ссылки на подписку
	uuid, err := getSubId(w, r)
	if err != nil {
		log.Print(err.Error(), " DeleteSub method")
		return
	}

	err = h.s.DeleteSub(uuid)
	if err != nil {
		http.Error(w, " ошибка при удалении записи о подписке", http.StatusInternalServerError)
		log.Print(err.Error(), " DeleteSub method")
		return
	}

	writeJSON(w, http.StatusOK, "задача успешно удалена")
}

// Метод для чтения записей о подписках

func (h *Handlers) ReadSubs(w http.ResponseWriter, r *http.Request) {

}

// Метод для показа функционала/методов для путей запроса

func (h *Handlers) ShowMethods(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) ShowSubscSum(w http.ResponseWriter, r *http.Request) {

}
