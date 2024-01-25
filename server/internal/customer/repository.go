package customer

import (
	"database/sql"
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
	// SeedData() error
}

type repository struct {
	db sql.DB
}

func NewCustomerRepository(db sql.DB) CustomerRepositoryInterface {

	return &repository{
		db: db,
	}
}
