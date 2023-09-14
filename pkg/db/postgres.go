package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "danu"
	password = "macrii"
	dbname   = "postgres"
)

var DB *gorm.DB

func InitPostgres() *gorm.DB {

	fmt.Println("Connecting to database...")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	DB, err := gorm.Open(postgres.Open(psqlconn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed postgres connection: %v\n", err)
	}

	return DB

}
