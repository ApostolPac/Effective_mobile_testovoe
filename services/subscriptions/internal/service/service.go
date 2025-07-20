package service

import (
	"subscriptions/internal/models"
	"time"
)

type Storage interface {
	CreateSubRequest(sub models.Subscription)(models.Subscription, error)
	ReadSubRequest(id int)(models.Subscription, error)
	ReadSubsRequest(userId string)([]models.Subscription, error)
	UpdateSubRequest(sub models.Subscription) error
	DeleteSubRequest(id int) error
	ShowSubscSumRequest(startPeriod time.Time, EndPeriod time.Time)([]models.Subscription, error)
}

type Service struct{
	s Storage
}

func NewService(a Storage) *Service{
	return &Service{
		s:a,
	}	
}

func CreateSub(sub models.Subscription) error {

	return nil
}

func ReadSub(id int) (*models.Subscription, error) {

	return nil
}

func ReadSubs(userId string) ([]models.Subscription, error) {

	return nil
}

func UpdateSub(sub models.Subscription) error {

	return nil
}

func DeleteSub(id int) error {

	return nil
}

func ShowSubscSum(startPeriod time.Time, EndPeriod time.Time) error {

	return nil
}
