package delivery

import (
	"database/sql"
)

// Delivery represents the structure of the deliveries data
type Delivery struct {
	ID                int
	OrderItemID       int
	DeliveredQuantity int
}

type DeliveryRepositoryInterface interface {
	// SeedData() error
}

type repository struct {
	db sql.DB
}

func NewDeliveryRepository(db sql.DB) DeliveryRepositoryInterface {
	return &repository{
		db: db,
	}
}
