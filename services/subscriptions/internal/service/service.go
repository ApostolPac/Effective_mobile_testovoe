package service

import (
	"log"
	"subscriptions/internal/models"
)


type Storage interface {
	CreateSubRequest(sub models.Subscription)(string, error)
	ReadSubRequest(id int) (*models.Subscription, error)
	ReadSubsRequest(userId string) ([]models.Subscription, error)
	UpdateSubRequest(sub models.Subscription) error
	DeleteSubRequest(id int) error
	ShowSubscSumRequest(serviceName string, userId string, startPeriod string, EndPeriod string) ([]models.Subscription, error)
}

type ServiceMethods struct {
	s Storage
}

func NewService(a Storage) *ServiceMethods {
	return &ServiceMethods{
		s: a,
	}
}

func (service *ServiceMethods) CreateSub(sub models.Subscription)(string, error) {
	userId, err := service.s.CreateSubRequest(sub)
	if err != nil {
		log.Print(err.Error(), "CreateSub method")
		return "", err
	}
	return userId, nil
}

func (service *ServiceMethods) ReadSub(id int) (*models.Subscription, error) {
	sub, err := service.s.ReadSubRequest(id)

	if err != nil {
		log.Printf("ReadSub method: error:%v", err.Error())
		return nil, err
	}

	return sub, nil
}

func (service *ServiceMethods) ReadSubs(userId string) ([]models.Subscription, error) {
	subs, err := service.s.ReadSubsRequest(userId)

	if err != nil {
		log.Printf("ReadSubs method: error:%v", err.Error())
		return nil, err
	}

	return subs, nil
}

func (service *ServiceMethods) UpdateSub(sub models.Subscription) error {
	err := service.s.UpdateSubRequest(sub)

	if err != nil {
		log.Printf("UpdateSub method: error:%v", err.Error())
		return err
	}

	return nil
}

func (service *ServiceMethods) DeleteSub(id int) error {
	err := service.s.DeleteSubRequest(id)

	if err != nil {
		log.Printf("DeleteSub method: error:%v", err.Error())
		return err
	}

	return nil
}

func (service *ServiceMethods) ShowSubscSum(serviceName string, userId string, startPeriod string, EndPeriod string) ([]models.Subscription, error) {

	subs, err := service.s.ShowSubscSumRequest(serviceName, userId, startPeriod, EndPeriod)

	if err != nil {
		log.Printf("ShowSubscSum method: error:%v", err.Error())
		return nil, err
	}

	return subs, nil
}
