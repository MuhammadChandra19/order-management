package order

//go:generate mockgen -source repository.go -destination mock/repository_mock.go -package=mock

import (
	"database/sql"
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
	Id           string    `json:"id"`
	OrderName    string    `json:"order_name"`
	CompanyName  string    `json:"company_name"`
	CustomerName string    `json:"customer_name"`
	ProductName  string    `json:"product_name"`
	OrderDate    time.Time `json:"order_date"`
	Delivered    string    `json:"delivered"`
	TotalAmount  float64   `json:"total_amount"`
}

type OrderRepositoryInterface interface {
	GetOrderList(search string, startDate, endDate time.Time, sortDirection string, limit, offset int) ([]OrderInfo, error)
}

type repository struct {
	db *sql.DB
}

func (r *repository) GetOrderList(search string, startDate, endDate time.Time, sortDirection string, limit, offset int) ([]OrderInfo, error) {
	args := []interface{}{}

	query := `
		SELECT
			o.id,
			o.order_name,
			cc.company_name,
			c.name as customer_name,
			oi.product as product_name,
			o.created_at as order_date,
			COALESCE(SUM(d.delivered_quantity), 0) as delivered,
			COALESCE(SUM(oi.price_per_unit * oi.quantity), 0) as total_amount
		FROM 
			orders o
		JOIN 
			customers c ON o.customer_id = c.user_id
		JOIN 
			customer_companies cc ON c.company_id = cc.company_id
		JOIN 
			order_items oi ON o.id = oi.order_id
		JOIN 
			deliveries d ON oi.id = d.order_item_id
		WHERE 1=1
	`

	// Add conditions for filtering
	if search != "" {
		paramName := strconv.Itoa(len(args) + 1)
		query += " AND (o.order_name ILIKE $" + paramName + " OR oi.product ILIKE $" + paramName + ")"
		args = append(args, "%"+search+"%")
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		query += " AND o.created_at BETWEEN $" + strconv.Itoa(len(args)+1) + "::timestamp AND $" + strconv.Itoa(len(args)+2) + "::timestamp"
		args = append(args, startDate, endDate)
	}

	orderByClause := "o.created_at DESC"

	if sortDirection == "ASC" {
		orderByClause = "o.created_at ASC"
	}
	query += `
		GROUP BY 
			o.id, cc.company_id, c.id, oi.product
		ORDER BY
		` + orderByClause + `
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
			&orderInfo.Id,
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

func NewOrderRepository(db *sql.DB) OrderRepositoryInterface {
	return &repository{
		db: db,
	}
}
