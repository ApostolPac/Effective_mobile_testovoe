package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"subscriptions/internal/models"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	createSub        = "INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING user_id"
	updateSub        = "UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6"
	deleteSub        = "DELETE FROM subscriptions WHERE id = $1"
	readSub          = "SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1"
	readSubs         = "SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE user_id = $1"
	showsubssum      = "SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE user_id = $1 AND service_name = $2 AND start_date >= $3 AND end_date   <= $4 ORDER BY id"
	showsubstotalsum = "SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE user_id = $1 AND service_name = $2 AND start_date >= $3 AND end_date   <= $4"
)

type Storage struct {
	Db *sql.DB
}

func NewStorage() *Storage {
	postUrl := os.Getenv("POSTGRES_DB_URL")
	db, err := sql.Open("postgres", postUrl)
	if err != nil {
		log.Fatalf("error during database start, %v", err)
		return nil
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatalf("error during connect to database, %v", err)
		return nil
	}

	return &Storage{Db: db}
}

func TimeForm(start time.Time, end time.Time) (string, string) {
	return start.Format("01-2006"), end.Format("01-2006")
}

func (s *Storage) RunMigrations() {
	driver, err := postgres.WithInstance(s.Db, &postgres.Config{})
	if err != nil {
		log.Fatalf("error during database driver initialization %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://internal/migrations", "postgres", driver)

	if err != nil {
		log.Fatalf("error during migration initialization: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("file application error migration: %v", err)
	}
}

func (s *Storage) CreateSubRequest(sub models.Subscription) (string, error) {
	var id string

	err := s.Db.QueryRow(createSub, sub.ServiceName, sub.Price, sub.UserId, sub.StartDate, sub.EndDate).Scan(&id)
	if err != nil {
		log.Printf("CreateSubRequest:error during creation of subscription record, error: %v", err.Error())
		return "", err
	}

	return id, nil
}

func (s *Storage) ReadSubRequest(id int) (*models.Subscription, error) {
	var sub models.Subscription
	var startDate, endDate time.Time

	err := s.Db.QueryRow(readSub, id).Scan(&sub.Id, &sub.ServiceName, &sub.Price, &sub.UserId, &startDate, &endDate)

	if err == sql.ErrNoRows {
		log.Printf("ReadSubRequest: error during read of subscription record, error: %v", err.Error())
		return nil, fmt.Errorf("no subscription record in database")
	}

	if err != nil {
		log.Printf("ReadSubRequest: error during read of subscription record %v", err.Error())
		return nil, err
	}

	sub.StartDate, sub.EndDate = TimeForm(startDate, endDate)

	return &sub, nil
}

func (s *Storage) ReadSubsRequest(userId string) ([]models.Subscription, error) {
	var subs []models.Subscription

	rows, err := s.Db.Query(readSubs, userId)
	if err != nil {
		log.Printf("ReadSubsRequest: error during read of subscriptions records, error: %v", err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var sub models.Subscription
		var startDate, endDate time.Time
		err = rows.Scan(&sub.Id, &sub.ServiceName, &sub.Price, &sub.UserId, &startDate, &endDate)
		if err != nil {
			log.Printf("ReadSubsRequest: error during rowscan, error: %v", err.Error())
			return nil, err
		}
		sub.StartDate, sub.EndDate = TimeForm(startDate, endDate)
		subs = append(subs, sub)
	}

	return subs, nil
}

func (s *Storage) UpdateSubRequest(sub models.Subscription) error {

	_, err := s.Db.Exec(updateSub, sub.ServiceName, sub.Price, sub.UserId, sub.StartDate, sub.EndDate, sub.Id)

	if err != nil {
		log.Printf("UpdateSubRequest: error during update of subscription record, error: %v", err.Error())
		return err
	}

	return nil
}

func (s *Storage) DeleteSubRequest(id int) error {

	_, err := s.Db.Exec(deleteSub, id)

	if err != nil {
		log.Printf("DeleteSubRequest: error during delete of subscription record, error: %v", err.Error())
		return err
	}

	return nil
}

func (s *Storage) ShowSubscSumRequest(serviceName string, userId string, startPeriod string, endPeriod string) ([]models.Subscription, error) {
	var subs []models.Subscription

	rows, err := s.Db.Query(showsubssum, userId, serviceName, startPeriod, endPeriod)
	if err != nil {
		log.Printf("ShowSubscSumRequest: error during read of subscriptions records, error: %v", err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var sub models.Subscription
		var startDate, endDate time.Time
		err = rows.Scan(&sub.Id, &sub.ServiceName, &sub.Price, &sub.UserId, &startDate, &endDate)

		if err != nil {
			log.Printf("ShowSubscSumRequest: error during rowscan, error: %v", err.Error())
			return nil, err
		}

		sub.StartDate, sub.EndDate = TimeForm(startDate, endDate)

		subs = append(subs, sub)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("ShowSubscSumRequest rows: %w", err)
	}

	var total int

	err = s.Db.QueryRow(showsubstotalsum, userId, serviceName, startPeriod, endPeriod).Scan(&total)
	if err != nil {
		log.Printf("ShowSubscSumRequest: error during read of sum, error: %v", err.Error())
		return nil, err
	}
	subs = append(subs, models.Subscription{
		ServiceName: "Итого",
		TotalSum:    total,
	})

	return subs, nil
}
