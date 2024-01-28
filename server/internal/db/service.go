package db

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	onceDB sync.Once
	db     *sql.DB
)

func setup() {
	var err error
	onceDB.Do(func() {
		db, err = sql.Open("postgres", "port=5432 user=your_username dbname=ordermanagement password=your_password sslmode=disable")
	})
	if err != nil {
		log.Fatal(err)
	}
}

func NewDatabase() *sql.DB {
	setup()
	return db
}
