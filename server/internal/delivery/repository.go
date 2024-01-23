package delivery

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
)

// Delivery represents the structure of the deliveries data
type Delivery struct {
	ID                int
	OrderItemID       int
	DeliveredQuantity int
}

type DeliveryRepositoryInterface interface {
	SeedData() error
}

type repository struct {
	db sql.DB
}

func (r *repository) SeedData() error {
	err := r.populateDeliveries()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) populateDeliveries() error {
	csvFile, err := os.Open("server/internal/db/deliveries.csv")
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
		orderItemID := record[1]
		deliveredQuantity := record[2]

		_, err = r.db.Exec(`
			INSERT INTO deliveries (order_item_id, delivered_quantity)
			VALUES ($1, $2)
		`, orderItemID, deliveredQuantity)
		if err != nil {
			log.Fatal(err, "deliveries")
		}
	}

	return nil
}

func NewDeliveryRepository(db sql.DB) DeliveryRepositoryInterface {
	return &repository{
		db: db,
	}
}
