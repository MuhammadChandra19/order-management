package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/MuhammadChandra19/order-management/internal/db"
	"github.com/MuhammadChandra19/order-management/internal/utils"
)

var database *sql.DB

func main() {
	database = db.NewDatabase()
	initTables()

	err := populateCustomerCompanies()
	if err != nil {
		log.Printf("Error populateCustomerCompanies %s\n", err)
	}
	fmt.Println("populateCustomerCompanies...")

	err = populateCustomers()
	if err != nil {
		log.Printf("Error populateCustomers %s\n", err)
	}
	fmt.Println("populateCustomers...")

	err = populateOrders()
	if err != nil {
		log.Printf("Error populateOrders %s\n", err)
	}
	fmt.Println("populateOrders...")

	err = populateOrderItems()
	if err != nil {
		log.Printf("Error populateOrderItems %s\n", err)
	}
	fmt.Println("populateOrderItems...")

	err = populateDeliveries()
	if err != nil {
		log.Printf("Error populateDeliveries %s\n", err)
	}
	fmt.Println("populateDeliveries...")
}

func initTables() {
	_, err := database.Exec(`
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

func populateCustomerCompanies() error {

	records, err := readData("customer_companies.csv")
	if err != nil {
		return err
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, record := range records {
		companyId := record[0]
		name := record[1]

		temp := len(valueStrings) * 2
		valueStrings = append(valueStrings, "($"+strconv.Itoa(temp+1)+", $"+strconv.Itoa(temp+2)+")")
		valueArgs = append(valueArgs, companyId)
		valueArgs = append(valueArgs, name)
	}

	stmt := fmt.Sprintf("INSERT INTO customer_companies (company_id, company_name) VALUES %s", strings.Join(valueStrings, ","))
	_, err = database.Exec(stmt, valueArgs...)
	if err != nil {
		return err
	}
	return nil
}

func populateCustomers() error {

	records, err := readData("customers.csv")
	if err != nil {
		return err
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, record := range records {
		userID := record[0]
		login := record[1]
		password := record[2]
		name := record[3]

		// Convert company_id to integer
		companyID, err := strconv.Atoi(record[4])
		if err != nil {
			return err
		}

		creditCards := record[5]

		temp := len(valueStrings) * 6
		valueStrings = append(
			valueStrings,
			"($"+strconv.Itoa(temp+1)+
				", $"+strconv.Itoa(temp+2)+
				", $"+strconv.Itoa(temp+3)+
				", $"+strconv.Itoa(temp+4)+
				", $"+strconv.Itoa(temp+5)+
				", $"+strconv.Itoa(temp+6)+
				")",
		)
		valueArgs = append(valueArgs, userID)
		valueArgs = append(valueArgs, login)
		valueArgs = append(valueArgs, password)
		valueArgs = append(valueArgs, name)
		valueArgs = append(valueArgs, companyID)
		valueArgs = append(valueArgs, creditCards)
	}

	stmt := fmt.Sprintf("INSERT INTO customers (user_id, login, password, name, company_id, credit_cards) VALUES %s", strings.Join(valueStrings, ","))
	_, err = database.Exec(stmt, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func populateOrders() error {
	records, err := readData("orders.csv")
	if err != nil {
		return err
	}
	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, record := range records {
		createdAt, err := utils.CompileDate(record[1])
		if err != nil {
			return err
		}

		orderName := record[2]
		customerID := record[3]

		temp := len(valueStrings) * 3
		valueStrings = append(
			valueStrings,
			"($"+strconv.Itoa(temp+1)+
				", $"+strconv.Itoa(temp+2)+
				", $"+strconv.Itoa(temp+3)+
				")",
		)
		valueArgs = append(valueArgs, createdAt)
		valueArgs = append(valueArgs, orderName)
		valueArgs = append(valueArgs, customerID)
	}

	stmt := fmt.Sprintf("INSERT INTO orders (created_at, order_name, customer_id) VALUES %s", strings.Join(valueStrings, ","))
	_, err = database.Exec(stmt, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func populateOrderItems() error {
	records, err := readData("order_items.csv")
	if err != nil {
		return err
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, record := range records {
		orderID := record[1]
		// Set default value if price_per_unit is empty
		pricePerUnitStr := record[2]
		var pricePerUnit float64
		if pricePerUnitStr == "" {
			pricePerUnit = 0.0
		} else {
			pricePerUnit, err = strconv.ParseFloat(pricePerUnitStr, 64)
			if err != nil {
				log.Printf("Error parsing price_per_unit for row %+v: %s\n", record, err)
				continue
			}
		}

		// Set default value if quantity is empty
		quantityStr := record[3]
		var quantity int
		if quantityStr == "" {
			quantity = 0
		} else {
			quantity, err = strconv.Atoi(quantityStr)
			if err != nil {
				log.Printf("Error parsing quantity for row %+v: %s\n", record, err)
				continue
			}
		}
		product := record[4]

		temp := len(valueStrings) * 4
		valueStrings = append(
			valueStrings,
			"($"+strconv.Itoa(temp+1)+
				", $"+strconv.Itoa(temp+2)+
				", $"+strconv.Itoa(temp+3)+
				", $"+strconv.Itoa(temp+4)+
				")",
		)
		valueArgs = append(valueArgs, orderID)
		valueArgs = append(valueArgs, pricePerUnit)
		valueArgs = append(valueArgs, quantity)
		valueArgs = append(valueArgs, product)
	}

	stmt := fmt.Sprintf("INSERT INTO order_items (order_id, price_per_unit, quantity, product) VALUES %s", strings.Join(valueStrings, ","))
	_, err = database.Exec(stmt, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func populateDeliveries() error {

	records, err := readData("deliveries.csv")
	if err != nil {
		return err
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, record := range records {
		orderItemID := record[1]
		deliveredQuantity := record[2]

		temp := len(valueStrings) * 2
		valueStrings = append(
			valueStrings,
			"($"+strconv.Itoa(temp+1)+
				", $"+strconv.Itoa(temp+2)+
				")",
		)
		valueArgs = append(valueArgs, orderItemID)
		valueArgs = append(valueArgs, deliveredQuantity)
	}

	stmt := fmt.Sprintf("INSERT INTO deliveries (order_item_id, delivered_quantity) VALUES %s", strings.Join(valueStrings, ","))
	_, err = database.Exec(stmt, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func readData(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	// skip first line
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
