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
	log.Print("writeJSON method: start of request")

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)

	log.Printf("writeJSON method: start of encoding - header:%v, body for encoding:%v", w.Header(), v)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Print(err.Error(), " writeJSON method")
		return err
	}

	log.Printf("writeJSON method: encoding complited")

	return nil
}

// Функция для получения id записи

func getSubId(w http.ResponseWriter, r *http.Request) (subIdi int, err error) {
	log.Printf("getSubId method: start of method, path: %v", r.URL.Path)

	log.Printf("getSubId method: start of TrimPrefix")

	id := strings.TrimPrefix(r.URL.Path, "/subscriptions/")

	log.Printf("getSubId method: TrimPrefix complited id = %v", id)

	if id == "" {
		log.Printf("getSubId method: error during id extraction, id is empty")
		return 0, fmt.Errorf("getSubId method: error during id extraction, id is empty")
	}

	log.Printf("getSubId method: start of convertion id type string to type int")

	subId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("getSubId method: id must be number, incorrect format, error:%v", err.Error())
		return 0, fmt.Errorf("getSubId method: id must be number, incorrect format")
	}

	log.Printf("getSubId method: successful convertion id to type int")

	log.Printf("getSubId method: successful request complited, subscription id = %v", subId)

	return subId, nil
}

func getService(r *http.Request) (serviceName string, err error) {
	log.Printf("getService method: start of method, path: %v", r.URL.Path)
	log.Printf("getService method: start of TrimPrefix")

	name := strings.TrimPrefix(r.URL.Path, "/subscriptions/sum/")

	log.Printf("getService method: TrimPrefix complited name = %v", name)

	if name == "" {
		log.Printf("getService method: error during service name extraction, service name is empty, PATH = %v", r.URL.Path)
		return "", fmt.Errorf("getService method: error during service name extraction, service name is empty, PATH = %v", r.URL.Path)
	}

	log.Printf("getService method: successful request complited, serviceName = %v", name)
	return name, nil
}

// Функция для получения uuid пользователя

func getUserUuid(r *http.Request) (uuid string, err error) {
	log.Printf("getUserUuid method: start of method, path: %v", r.URL.Path)
	log.Printf("getUserUuid method: start of header extraction, path: %v", r.URL.Path)

	userUuid := r.Header.Get("Authorization")

	log.Printf("getUserUuid method: extraction of userUuid complited, userUuid = %v", userUuid)

	if userUuid == "" {
		log.Print("getUserUuid method: error during userUuid extraction, empty userUuid in header Authorization")
		return "", fmt.Errorf("getUserUuid method: error during userUuid extraction, empty userUuid in header Authorization")
	}

	log.Printf("getUserUuid method: successful request complited, userUuid = %v", userUuid)

	return userUuid, nil
}

// CreateSub godoc
// @Summary     Создать подписку
// @Description Создаёт новую подписку. Все данные, включая user_id, передаются в теле запроса.
// @Tags        subscriptions
// @Accept      json
// @Produce     json
// @Param       subscription  body   models.Subscription true "Данные новой подписки"
// @Success     200           {string} string           "Подписка создана"
// @Failure     400           {object} map[string]string "Bad Request"
// @Failure     500           {object} map[string]string "Internal Server Error"
// @Router      /subscriptions [post]
func (h *Handlers) CreateSub(w http.ResponseWriter, r *http.Request) {

	log.Printf("CreateSub: method=%s url=%s", r.Method, r.URL.Path)

	var sub models.Subscription

	log.Printf("CreateSub: start of sub decoding")

	err := json.NewDecoder(r.Body).Decode(&sub)

	log.Printf("CreateSub: sub decoding complited, sub: %v", sub)
	if err != nil {
		http.Error(w, dataStructError, http.StatusBadRequest)
		log.Print("CreateSub method: error durind decoding of json body ", err.Error())
		return
	}

	log.Printf("CreateSub: created sub %v ", sub)

	userId, err := h.w.AsyncCreateSub(sub)
	if err != nil {
		http.Error(w, "error during creation of a subscription record", http.StatusInternalServerError)
		log.Print("CreateSub method: error during AsyncCreateSub request ", err.Error())
		return
	}

	log.Printf("CreateSub method: creation of subscription record complited - userId: %v ", userId)

	answer := fmt.Sprintf("the subscription record has been successfully registered, user id = %v", userId)

	log.Printf("CreateSub method: start of request to writeJSON, answer = %v", answer)

	err = writeJSON(w, http.StatusOK, answer)
	if err != nil {
		http.Error(w, "error during writing answer", http.StatusInternalServerError)
		log.Print("CreateSub method: error during writeJSON ", err.Error())
		return
	}

	log.Print("CreateSub method: successful request complited")
}

// Хендлер для чтения записи о подписке

// ReadSub godoc
// @Summary     Получить подписку по ID
// @Description Возвращает данные одной подписки по её ID.
// @Tags        subscriptions
// @Produce     json
// @Param       id   path      int    true  "Subscription ID"
// @Success     200  {object}  models.Subscription "Данные подписки"
// @Failure     400  {object}  map[string]string   "Bad Request"
// @Failure     404  {object}  map[string]string   "Not Found"
// @Failure     500  {object}  map[string]string   "Internal Server Error"
// @Router      /subscriptions/{id} [get]
func (h *Handlers) ReadSub(w http.ResponseWriter, r *http.Request) {

	log.Printf("ReadSub: method=%v url=%v", r.Method, r.URL.Path)

	log.Printf("ReadSub: start of request to getSubId")

	id, err := getSubId(w, r)

	log.Printf("ReadSub: request to getSubId method complited, id = %v", id)
	if err != nil {
		http.Error(w, missedLinkError, http.StatusBadRequest)
		log.Print("ReadSub method: error during getSubId request ", err.Error())
		return
	}
	log.Printf("ReadSub: request to AsyncReadSub method, id = %v", id)

	sub, err := h.w.AsyncReadSub(models.Subscription{Id: id})

	log.Printf("ReadSub: request to AsyncReadSub method complited, sub = %v", sub)

	if err != nil {
		http.Error(w, "error during read of subscription record", http.StatusInternalServerError)
		log.Print("ReadSub method: error during AsyncReadSub request ", err.Error())
		return
	}

	log.Printf("ReadSub method: start of request to writeJSON, sub = %v", sub)

	err = writeJSON(w, http.StatusOK, sub)

	if err != nil {
		http.Error(w, "error during writing answer", http.StatusInternalServerError)
		log.Print("ReadSub method: error during writeJSON ", err.Error())
		return
	}

	log.Print("ReadSub method: successful request complited")
}


// UpdateSub godoc
// @Summary     Обновить подписку
// @Description Обновляет запись подписки: указывается ID в пути и новые данные в теле запроса.
// @Tags        subscriptions
// @Accept      json
// @Produce     json
// @Param       id            path   int                 true  "Subscription ID"
// @Param       subscription  body   models.Subscription true  "Новые данные подписки"
// @Success     200           {string} string           "Подписка обновлена"
// @Failure     400           {object} map[string]string "Bad Request"
// @Failure     500           {object} map[string]string "Internal Server Error"
// @Router      /subscriptions/{id} [put]
func (h *Handlers) UpdateSub(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateSub: method=%v url=%v", r.Method, r.URL.Path)

	var sub models.Subscription
	log.Printf("UpdateSub: start of sub decoding")

	err := json.NewDecoder(r.Body).Decode(&sub)
	log.Printf("UpdateSub: sub decoding complited, sub: %v", sub)

	if err != nil {
		http.Error(w, dataStructError, http.StatusBadRequest)
		log.Print("UpdateSub method: error durind decoding of json body ", err.Error())
		return
	}

	log.Printf("UpdateSub: start of request to getSubId")

	id, err := getSubId(w, r)

	log.Printf("UpdateSub: request to getSubId method complited, id = %v", id)
	if err != nil {
		http.Error(w, missedLinkError, http.StatusBadRequest)
		log.Print("UpdateSub method: error during getSubId request ", err.Error())
		return
	}

	sub.Id = id

	log.Printf("UpdateSub: request to AsyncUpdateSub method, id = %v", id)

	err = h.w.AsyncUpdateSub(sub)

	log.Printf("UpdateSub: request to AsyncUpdateSub method complited")

	if err != nil {
		http.Error(w, "error during subscription record update", http.StatusInternalServerError)
		log.Print("UpdateSub method: error during AsyncUpdateSub request ", err.Error())
		return
	}

	log.Printf("UpdateSub method: start of request to writeJSON")

	err = writeJSON(w, http.StatusOK, "subscrription record updated successfuly")
	if err != nil {
		http.Error(w, "error during writing answer", http.StatusInternalServerError)
		log.Print("UpdateSub method: error during writeJSON ", err.Error())
		return
	}

	log.Print("UpdateSub method: successful request complited")
}

// Хендлер для удаления записи о подписке

// DeleteSub godoc
// @Summary     Удалить подписку
// @Description Удаляет запись подписки по её ID.
// @Tags        subscriptions
// @Produce     json
// @Param       id   path      int    true  "Subscription ID"
// @Success     200  {string}  string "Подписка удалена"
// @Failure     400  {object}  map[string]string "Bad Request"
// @Failure     500  {object}  map[string]string "Internal Server Error"
// @Router      /subscriptions/{id} [delete]
func (h *Handlers) DeleteSub(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteSub: method=%v url=%v", r.Method, r.URL.Path)

	//удаляем часть пути для получения только ссылки на подписку

	log.Printf("DeleteSub: start of request to getSubId")

	id, err := getSubId(w, r)

	log.Printf("DeleteSub: request to getUserUuid method complited, id = %v", id)

	if err != nil {
		http.Error(w, missedLinkError, http.StatusBadRequest)
		log.Print(err.Error(), " DeleteSub method")
		return
	}

	log.Printf("DeleteSub: request to AsyncDeleteSub method, id = %v", id)

	err = h.w.AsyncDeleteSub(models.Subscription{Id: id})

	log.Printf("DeleteSub: request to AsyncDeleteSub method complited")

	if err != nil {
		http.Error(w, "error during deletion of subscription record", http.StatusInternalServerError)
		log.Printf("DeleteSub: error during request to AsyncDeleteSub, error = %v", err)
		return
	}

	log.Printf("DeleteSub method: start of request to writeJSON")

	err = writeJSON(w, http.StatusOK, "subscription record deleted successfuly")
	if err != nil {
		http.Error(w, "error during writing answer", http.StatusInternalServerError)
		log.Print("DeleteSub method: error during writeJSON ", err.Error())
		return
	}

	log.Print("DeleteSub method: successful request complited")
}

// Метод для чтения записей о подписках с общей суммой (показывает все подписки со всеми сервисами для конкретного пользователя)

// ReadSubs godoc
// @Summary     Получить все подписки пользователя
// @Description Возвращает список подписок для пользователя, UUID берётся из заголовка Authorization.
// @Tags        subscriptions
// @Produce     json
// @Param       Authorization header string true "User UUID"
// @Success     200 {array} models.Subscription "Список подписок"
// @Failure     400 {object} map[string]string      "Bad Request"
// @Failure     500 {object} map[string]string      "Internal Server Error"
// @Router      /subscriptions [get]
func (h *Handlers) ReadSubs(w http.ResponseWriter, r *http.Request) {
	log.Printf("ReadSubs: method=%v url=%v", r.Method, r.URL.Path)

	log.Printf("ReadSubs: start of request to getUserUuid")

	uuid, err := getUserUuid(r)

	log.Printf("ReadSubs: request to getUserUuid method complited, uuid = %v", uuid)

	if err != nil {
		http.Error(w, missedLinkError, http.StatusInternalServerError)
		log.Printf("ReadSubs: error during request to getUserUuid method, uuid = %v, error = %v", uuid, err.Error())
		return
	}

	log.Printf("ReadSubs: request to AsyncReadSubs method, uuid = %v", uuid)

	subs, err := h.w.AsyncReadSubs(models.Subscription{UserId: uuid})

	log.Printf("ReadSubs: request to AsyncReadSubs method complited, subs = %v", subs)

	if err != nil {
		http.Error(w, "error during subs extraction", http.StatusInternalServerError)
		log.Printf("ReadSubs: error during request to AsyncReadSubs, error = %v", err)
		return
	}

	log.Printf("ReadSubs method: start of request to writeJSON, subs = %v", subs)

	err = writeJSON(w, http.StatusOK, subs)
	if err != nil {
		http.Error(w, "error during writing answer", http.StatusInternalServerError)
		log.Print("ReadSubs method: error during writeJSON ", err.Error())
		return
	}

	log.Print("ReadSubs method: successful request complited")
}


// ShowSubscSum godoc
// @Summary     Получить подписки и их сумму по сервису за период
// @Description Сервис указывается в пути, период (start_date и end_date) — в теле, user UUID — в заголовке Authorization.
// @Tags        subscriptions
// @Accept      json
// @Produce     json
// @Param       Authorization header string               true  "User UUID"
// @Param       service       path   string               true  "Service name (например, Netflix)"
// @Param       period        body   models.ShowSubscSum  true  "Период в формате для примера {2025-08-01T00:00:00Z}"
// @Success     200           {array}  models.Subscription
// @Failure     400           {object} map[string]string  "Bad Request"
// @Failure     500           {object} map[string]string  "Internal Server Error"
// @Router      /subscriptions/sum/{service} [post]
func (h *Handlers) ShowSubscSum(w http.ResponseWriter, r *http.Request) {
	log.Printf("ShowSubscSum: method=%v url=%v", r.Method, r.URL.Path)
	log.Printf("ShowSubscSum: request to getService method")

	serviceName, err := getService(r)

	log.Printf("ShowSubscSum: request to getService method complited")

	if err != nil {
		http.Error(w, missedLinkError, http.StatusBadRequest)
		log.Printf("ShowSubscSum: error durind request to getService, error: %v", err)
		return
	}

	log.Printf("ShowSubscSum: successful request to getService method, serviceName = %v", serviceName)

	log.Printf("ShowSubscSum: start of request to getUserUuid")

	uuid, err := getUserUuid(r)

	log.Printf("ShowSubscSum: request to getUserUuid method complited, uuid = %v", uuid)

	if err != nil {
		http.Error(w, missedLinkError, http.StatusInternalServerError)
		log.Printf("ShowSubscSum: error during request to getUserUuid method, uuid = %v, error = %v", uuid, err.Error())
		return
	}

	var periods models.ShowSubscSum

	log.Printf("ShowSubscSum: start of decoding periods")

	err = json.NewDecoder(r.Body).Decode(&periods)

	log.Printf("ShowSubscSum: decoding complited successfuly, periods = %v", periods)

	if err != nil {
		http.Error(w, dataStructError, http.StatusBadRequest)
		log.Printf("ShowSubscSum: error during periods decoding, periods = %v, error = %v", periods, err.Error())
		return
	}

	log.Printf("ShowSubscSum: decoding complited successfuly, periods = %v", periods)

	log.Printf("ShowSubscSum: start of request to AsyncShowSubscSum, ServiceName = %v, UserId = %v, StartDate = %v, EndDate = %v", serviceName, uuid, periods.StartDate, periods.EndDate)

	subs, err := h.w.AsyncShowSubscSum(models.Subscription{ServiceName: serviceName, UserId: uuid, StartDate: periods.StartDate, EndDate: periods.EndDate})

	log.Printf("ShowSubscSum: complited request to AsyncShowSubscSum, subs = %v", subs)

	if err != nil {
		http.Error(w, "error during subscription records summation", http.StatusInternalServerError)
		log.Printf("ShowSubscSum: error during request to AsyncShowSubscSum, error: %v", err)
		return
	}

	log.Printf("ShowSubscSum method: start of request to writeJSON, subs = %v", subs)

	err = writeJSON(w, http.StatusOK, subs)
	if err != nil {
		log.Print("CreateSub method: error during writeJSON ", err.Error())
		http.Error(w, "error during writing answer", http.StatusInternalServerError)
		return
	}

	log.Print("ShowSubscSum method: successful request complited")
}
