package service

import (
	"fmt"
	"log"
	"subscriptions/internal/models"
)

type JobType string

const (
	JobCreate  JobType = "create"
	JobUpdate  JobType = "update"
	JobDelete  JobType = "delete"
	JobShowOne JobType = "show_one"
	JobShowAll JobType = "show_all"
	JobShowSum JobType = "show_all_sum"
)

type Service interface {
	CreateSub(sub models.Subscription) (string, error)                                                                   // Метод для создания записи. Возвращает id пользователя и ошибку.
	ReadSub(id int) (*models.Subscription, error)                                                                        // Метод для чтения записи по её id.
	ReadSubs(userId string) ([]models.Subscription, error)                                                               // Метод для чтения среза записей для конкретного пользователя.
	UpdateSub(sub models.Subscription) error                                                                             // Метод для обновления записей методом Update.
	DeleteSub(id int) error                                                                                              // Метод для удаления записи о подписке.
	ShowSubscSum(serviceName string, userId string, startPeriod string, EndPeriod string) ([]models.Subscription, error) // Метод для получения сум подписок, для начала работы нужно -
	// отправить период внутри которого будем искать записи о подписках
}

type Job struct {
	Type    JobType
	Request models.Subscription
	Result  chan JobResult
}

type JobResult struct {
	Result interface{}
	Error  error
}

type WorkerPool struct {
	s Service
}

var jobChan chan Job

func StartWorkerPool(numWorkers int, se Service) *WorkerPool {
	jobChan = make(chan Job, 1000)

	wp := WorkerPool{s: se}

	for i := 0; i < numWorkers; i++ {
		go wp.worker(i, jobChan)
	}

	return &wp
}

func (w *WorkerPool) worker(i int, jobs <-chan Job) {

	var err error
	var result interface{}
	for job := range jobs {
		log.Printf("горутина %v получила задачу", i)
		switch job.Type {
		case JobCreate:
			result, err = w.s.CreateSub(job.Request)
		case JobUpdate:
			err = w.s.UpdateSub(job.Request)
		case JobDelete:
			err = w.s.DeleteSub(job.Request.Id)
		case JobShowOne:
			result, err = w.s.ReadSub(job.Request.Id)
		case JobShowAll:
			result, err = w.s.ReadSubs(job.Request.UserId)
		case JobShowSum:
			result, err = w.s.ShowSubscSum(job.Request.ServiceName, job.Request.UserId, job.Request.StartDate, job.Request.EndDate)
		}
		log.Printf("горутина %v выполнила задачу", i)
		job.Result <- JobResult{Result: result, Error: err}
	}

}

func (w *WorkerPool) AsyncCreateSub(sub models.Subscription) (string, error) {
	jobresult := make(chan JobResult, 1)

	jobChan <- Job{Type: JobCreate, Request: sub, Result: jobresult}

	res := <-jobresult

	strResult, ok := res.Result.(string)

	if !ok {
		return "", res.Error
	}

	return strResult, res.Error
}

func (w *WorkerPool) AsyncUpdateSub(sub models.Subscription) error {
	jobresult := make(chan JobResult, 1)

	jobChan <- Job{Type: JobUpdate, Request: sub, Result: jobresult}

	res := <-jobresult

	return res.Error
}

func (w *WorkerPool) AsyncDeleteSub(sub models.Subscription) error {
	jobresult := make(chan JobResult, 1)

	jobChan <- Job{Type: JobDelete, Request: sub, Result: jobresult}

	res := <-jobresult

	return res.Error
}

func (w *WorkerPool) AsyncReadSub(sub models.Subscription) (*models.Subscription, error) {
	jobresult := make(chan JobResult, 1)

	jobChan <- Job{Type: JobShowOne, Request: sub, Result: jobresult}

	res := <-jobresult

	subscr, ok := res.Result.(*models.Subscription)

	if !ok || subscr == nil {
		return nil, fmt.Errorf("неподходящий тип или отстутствует подписка, %v", ok)
	}

	return subscr, res.Error
}

func (w *WorkerPool) AsyncReadSubs(sub models.Subscription) ([]models.Subscription, error) {
	jobresult := make(chan JobResult, 1)

	jobChan <- Job{Type: JobShowAll, Request: sub, Result: jobresult}

	res := <-jobresult

	subscriptions, ok := res.Result.([]models.Subscription)

	if !ok || subscriptions == nil {
		return nil, fmt.Errorf("неподходящий тип или отстутствуют подписки, %v", ok)
	}

	return subscriptions, res.Error
}

func (w *WorkerPool) AsyncShowSubscSum(sub models.Subscription) ([]models.Subscription, error) {
	jobresult := make(chan JobResult, 1)

	jobChan <- Job{Type: JobShowSum, Request: sub, Result: jobresult}

	res := <-jobresult

	subscriptions, ok := res.Result.([]models.Subscription)

	if !ok || subscriptions == nil {
		return nil, fmt.Errorf("неподходящий тип или отстутствуют подписки, %v", ok)
	}

	return subscriptions, res.Error
}
