package customer

import (
	"database/sql"
	"encoding/csv"
	"os"
	"strconv"
)

// CustomerCompany represents the structure of the customer_companies data
type CustomerCompany struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Customer represents the structure of the customers data
type Customer struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	CompanyID   int    `json:"company_id"`
	CreditCards string `json:"credit_cards"`
}

type CustomerRepositoryInterface interface {
	SeedData() error
}

type repository struct {
	db sql.DB
}

func (r *repository) SeedData() error {
	err := r.populateCustomerCompanies()
	if err != nil {
		return err
	}

	err = r.populateCustomers()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) populateCustomerCompanies() error {
	csvFile, err := os.Open("server/internal/db/customer_companies.csv")
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records[1:] {
		companyId := record[0]
		name := record[1]
		_, err := r.db.Exec("INSERT INTO customer_companies (company_id, company_name) VALUES ($1, $2)", companyId, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) populateCustomers() error {
	csvFile, err := os.Open("server/internal/db/customers.csv")
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records[1:] {
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

		_, err = r.db.Exec(`
			INSERT INTO customers (user_id, login, password, name, company_id, credit_cards)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, userID, login, password, name, companyID, creditCards)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewCustomerRepository(db sql.DB) CustomerRepositoryInterface {

	return &repository{
		db: db,
	}
}
