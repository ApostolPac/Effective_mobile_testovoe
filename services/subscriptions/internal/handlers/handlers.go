package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"subscriptions/internal/models"
)

const (
	dataStructError = "неправильный формат данных"
	missedLinkError = "ссылка на запись о подписке отсутствует"
)

type WorkerPool interface {
	AsyncCreateSub(sub models.Subscription) (string, error)
	AsyncUpdateSub(sub models.Subscription) error
	AsyncDeleteSub(sub models.Subscription) error
	AsyncReadSub(sub models.Subscription) (*models.Subscription, error)
	AsyncReadSubs(sub models.Subscription) ([]models.Subscription, error)
	AsyncShowSubscSum(sub models.Subscription) ([]models.Subscription, error)
}

type Handlers struct {
	w WorkerPool
}

func NewHandler(a WorkerPool) *Handlers {
	return &Handlers{
		w: a,
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

func getService(w http.ResponseWriter, r *http.Request) (serviceName string, err error) {
	name := strings.TrimPrefix(r.URL.Path, "/subscriptions/sum/")

	if name == "" {
		log.Print("ошибка при извлечении service_name - отсутствует имя сервиса")
		return "", fmt.Errorf("service_name подписок отсутствует")
	}

	return name, nil
}

// Функция для получения uuid пользователя

func getUserUuid(w http.ResponseWriter, r *http.Request) (uuid string, err error) {

	userUuid := r.Header.Get("Authorization")

	if userUuid == "" {
		log.Print("empty user_id in header Authorization", " getUserUuid method")
		return "", fmt.Errorf("отстуствует header с user_id")
	}

	return userUuid, nil
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

	userId, err := h.w.AsyncCreateSub(sub)
	if err != nil {
		http.Error(w, "ошибка при создании записи о подписке", http.StatusInternalServerError)
		log.Print(err.Error(), " CreateSub method")
		return
	}

	answer := fmt.Sprintf("запись о подписке успешно зарегистрирована, id записи = %v, id пользователя = %v", sub.Id, userId)

	err = writeJSON(w, http.StatusOK, answer)
	if err != nil {
		log.Print(err.Error(), " CreateSub method")
		return
	}

}

// Хендлер для чтения записи о подписке

func (h *Handlers) ReadSub(w http.ResponseWriter, r *http.Request) {
	id, err := getSubId(w, r)
	if err != nil {
		log.Print(err.Error(), " ReadSub method")
		return
	}

	sub, err := h.w.AsyncReadSub(models.Subscription{Id: id})

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

	id, err := getSubId(w, r)
	if err != nil {
		log.Print(err.Error(), " ReadSub method")
		return
	}

	sub.Id = id
	
	err = h.w.AsyncUpdateSub(sub)
	if err != nil {
		http.Error(w, "ошибка при обновлении записи о подписке", http.StatusInternalServerError)
		log.Print(err.Error(), " UpdateSub method")
		return
	}

	writeJSON(w, http.StatusOK, "запись о подписке успешно изменена")
}

// Хендлер для удаления записи о подписке

func (h *Handlers) DeleteSub(w http.ResponseWriter, r *http.Request) {
	//удаляем часть пути для получения только ссылки на подписку
	id, err := getSubId(w, r)
	if err != nil {
		log.Print(err.Error(), " DeleteSub method")
		return
	}

	err = h.w.AsyncDeleteSub(models.Subscription{Id: id})
	if err != nil {
		http.Error(w, " ошибка при удалении записи о подписке", http.StatusInternalServerError)
		log.Print(err.Error(), " DeleteSub method")
		return
	}

	writeJSON(w, http.StatusOK, "запись о подписке успешно удалена")
}

// Метод для чтения записей о подписках с общей суммой (показывает все подписки со всеми сервисами для конкретного пользователя)

func (h *Handlers) ReadSubs(w http.ResponseWriter, r *http.Request) {

	uuid, err := getUserUuid(w, r)
	if err != nil {
		http.Error(w, " ошибка при чтении записей о подписках", http.StatusInternalServerError)
		log.Print(err.Error(), " ReadSubs method")
		return
	}

	subs, err := h.w.AsyncReadSubs(models.Subscription{UserId: uuid})

	if err != nil {
		http.Error(w, " ошибка при чтении записей о подписках", http.StatusInternalServerError)
		log.Print(err.Error(), " ReadSubs method")
		return
	}

	writeJSON(w, http.StatusOK, subs)

}

// Метод для получения записей о подписке с общей суммой по конкретному

func (h *Handlers) ShowSubscSum(w http.ResponseWriter, r *http.Request) {
	serviceName, err := getService(w, r)

	if err != nil {
		http.Error(w, missedLinkError, http.StatusBadRequest)
		return
	}

	uuid, err := getUserUuid(w, r)

	if err != nil {
		http.Error(w, " ошибка при чтении записей о подписках", http.StatusInternalServerError)
		log.Print(err.Error(), " ReadSubs method")
		return
	}

	var periods models.ShowSubscSum

	err = json.NewDecoder(r.Body).Decode(&periods)
	if err != nil {
		http.Error(w, dataStructError, http.StatusBadRequest)
		log.Print(err.Error())
		return
	}

	subs, err := h.w.AsyncShowSubscSum(models.Subscription{ServiceName: serviceName, UserId:uuid, StartDate: periods.StartDate, EndDate: periods.EndDate})

	if err != nil {
		http.Error(w, " ошибка при подсчёте суммы подписок", http.StatusInternalServerError)
		log.Print(err.Error(), " ShowSubscSum method")
		return
	}

	writeJSON(w, http.StatusOK, subs)
}
