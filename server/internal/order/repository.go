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

// OrderRepositoryInterface defines the interface for accessing order-related data.

type OrderRepositoryInterface interface {
	// GetOrderList retrieves a list of orders based on search criteria, date range, sorting, limit, and offset.
	//
	// Parameters:
	// - search (string): Search string for filtering orders by name or product.
	//
	// - startDate (time.Time): Start date for date range filtering.
	//
	// - endDate (time.Time): End date for date range filtering.
	//
	// - sortDirection (string): Sorting direction (ASC or DESC).
	//
	// - limit (int): Maximum number of orders to retrieve.
	//
	// - offset (int): Pagination offset.
	//
	// Returns:
	// - ([]*OrderInfo, error): Retrieved orders and any retrieval errors.
	GetOrderList(search string, startDate, endDate time.Time, sortDirection string, limit, offset int) ([]*OrderInfo, error)
}

type repository struct {
	db *sql.DB
}

func (r *repository) GetOrderList(search string, startDate, endDate time.Time, sortDirection string, limit, offset int) ([]*OrderInfo, error) {
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

	var ordersInfo []*OrderInfo

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
		ordersInfo = append(ordersInfo, &orderInfo)
	}

	return ordersInfo, nil
}

// NewOrderRepository creates a new instance of the OrderRepositoryInterface,
// which is an interface for interacting with the database to perform operations
// related to orders.
//
// Parameters:
// - db (*sql.DB): A pointer to the database instance to be used by the repository.
//
// Returns:
// - OrderRepositoryInterface: An instance of the OrderRepositoryInterface.
//
// Notes:
//   - The function creates a new repository struct with the provided database connection.
//   - The returned repository implements the OrderRepositoryInterface.
//   - The OrderRepositoryInterface defines methods for interacting with order-related data
//     in the database.
func NewOrderRepository(db *sql.DB) OrderRepositoryInterface {
	return &repository{
		db: db,
	}
}
