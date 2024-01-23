package db

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

// type database struct {
// 	db          *sql.DB
// 	isPopulated bool
// }

var (
	onceDB sync.Once
	db     *sql.DB
)

func setup() bool {
	var err error
	onceDB.Do(func() {
		db, err = sql.Open("postgres", "port=5432 user=your_username dbname=ordermanagement password=your_password sslmode=disable")
	})
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Initialize tables
	initTables()
	// Check if data has been populated
	isPopulated := isDataPopulated()
	if !isDataPopulated() {

		// Set the data population flag to true
		setDataPopulatedFlag(true)

	}

	return isPopulated
}

func NewDatabase() (sql.DB, bool) {
	isPopulated := setup()
	return *db, isPopulated
}

func initTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS customer_companies (
			company_id SERIAL PRIMARY KEY,
			company_name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			user_id TEXT UNIQUE NOT NULL,
			login TEXT NOT NULL,
			password TEXT NOT NULL,
			name TEXT NOT NULL,
			company_id INT REFERENCES customer_companies(company_id),
			credit_cards TEXT
		);

		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP,
			order_name TEXT NOT NULL,
			customer_id TEXT REFERENCES customers(user_id)
		);

		CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INT REFERENCES orders(id),
			price_per_unit NUMERIC,
			quantity INT NOT NULL,
			product TEXT
		);

		CREATE TABLE IF NOT EXISTS deliveries (
			id SERIAL PRIMARY KEY,
			order_item_id INT REFERENCES order_items(id),
			delivered_quantity INT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS data_population_flag (
			id SERIAL PRIMARY KEY,
			populated BOOLEAN NOT NULL
		);
		
		INSERT INTO data_population_flag (populated) VALUES (false);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func isDataPopulated() bool {
	var populated bool
	err := db.QueryRow("SELECT populated FROM data_population_flag WHERE id = 1").Scan(&populated)
	if err != nil {
		log.Fatal(err)
	}
	return populated
}

func setDataPopulatedFlag(populated bool) {
	_, err := db.Exec("UPDATE data_population_flag SET populated = $1 WHERE id = 1", populated)
	if err != nil {
		log.Fatal(err)
	}
}
