package storage

import(
	"log"
	"os"
	"github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file")


type Storage struct{

}

func NewStorage()*Storage{
	return &Storage{

	}
}
