package order

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

// OrderItem represents the structure of the order_items data
type OrderItem struct {
	ID           int
	OrderID      int
	PricePerUnit float64
	Quantity     int
	Product      string
}

// Order represents the structure of the orders data
type Order struct {
	ID         int
	CreatedAt  time.Time
	OrderName  string
	CustomerID string
}

// OrderInfo represents the structure of the information to be shown
type OrderInfo struct {
	OrderName    string    `json:"order_name"`
	CompanyName  string    `json:"company_name"`
	CustomerName string    `json:"customer_name"`
	ProductName  string    `json:"product_name"`
	OrderDate    time.Time `json:"order_date"`
	Delivered    string    `json:"delivered"`
	TotalAmount  float64   `json:"total_amount"`
}

type OrderRepositoryInterface interface {
	SeedData() error
	GetOrderList(search string, startDate, endDate time.Time, limit, offset int) ([]OrderInfo, error)
}

type repository struct {
	db sql.DB
}

func (r *repository) GetOrderList(search string, startDate, endDate time.Time, limit, offset int) ([]OrderInfo, error) {
	args := []interface{}{}

	query := `
		SELECT
			o.order_name,
			cc.company_name,
			c.name as customer_name,
			oi.product as product_name,
			o.created_at as order_date,
			COALESCE(SUM(d.delivered_quantity), 0) as delivered,
			COALESCE(SUM(oi.price_per_unit * oi.quantity), 0) as total_amount
		FROM orders o
		JOIN customers c ON o.customer_id = c.user_id
		JOIN customer_companies cc ON c.company_id = cc.company_id
		JOIN order_items oi ON o.id = oi.order_id
		JOIN deliveries d ON oi.id = d.order_item_id
		WHERE 1=1
	`

	// Add conditions for filtering
	if search != "" {
		paramName := strconv.Itoa(len(args) + 1)
		query += " AND (o.order_name ILIKE $" + paramName + " OR oi.product ILIKE $" + paramName + ")"
		args = append(args, "%"+search+"%")
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		query += " AND o.created_at >= $" + strconv.Itoa(len(args)+1) + "AND o.created_at <= $" + strconv.Itoa(len(args)+1)
		args = append(args, startDate, endDate)
	} else if !startDate.IsZero() {
		query += " AND o.created_at >= $" + strconv.Itoa(len(args)+1)
		args = append(args, startDate)
	} else if !endDate.IsZero() {
		query += " AND o.created_at <= $" + strconv.Itoa(len(args)+1)
		args = append(args, endDate)
	}

	query += `
		GROUP BY o.id, cc.company_id, c.id, oi.product
		ORDER BY o.created_at DESC
		LIMIT $` + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordersInfo []OrderInfo

	for rows.Next() {
		var orderInfo OrderInfo
		err := rows.Scan(
			&orderInfo.OrderName,
			&orderInfo.CompanyName,
			&orderInfo.CustomerName,
			&orderInfo.ProductName,
			&orderInfo.OrderDate,
			&orderInfo.Delivered,
			&orderInfo.TotalAmount,
		)
		if err != nil {
			return nil, err
		}
		ordersInfo = append(ordersInfo, orderInfo)
	}

	return ordersInfo, nil
}

func (r *repository) SeedData() error {
	err := r.populateOrders()
	if err != nil {
		return err
	}

	err = r.populateOrderItems()
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) populateOrders() error {
	csvFile, err := os.Open("server/internal/db/orders.csv")
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
		createdAt, err := time.Parse(time.RFC3339, record[1])
		if err != nil {
			return err
		}

		orderName := record[2]
		customerID := record[3]

		_, err = r.db.Exec(`
			INSERT INTO orders (created_at, order_name, customer_id)
			VALUES ($1, $2, $3)
		`, createdAt, orderName, customerID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) populateOrderItems() error {
	csvFile, err := os.Open("server/internal/db/order_items.csv")
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
		orderID := record[1]
		// Set default value if price_per_unit is empty
		pricePerUnitStr := record[2]
		var pricePerUnit float64
		if pricePerUnitStr == "" {
			pricePerUnit = 0.0 // You can set any default value here
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
			quantity = 0 // You can set any default value here
		} else {
			quantity, err = strconv.Atoi(quantityStr)
			if err != nil {
				log.Printf("Error parsing quantity for row %+v: %s\n", record, err)
				continue
			}
		}
		product := record[4]

		_, err = r.db.Exec(`
			INSERT INTO order_items (order_id, price_per_unit, quantity, product)
			VALUES ($1, $2, $3, $4)
		`, orderID, pricePerUnit, quantity, product)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewOrderRepository(db sql.DB) OrderRepositoryInterface {
	return &repository{
		db: db,
	}
}
