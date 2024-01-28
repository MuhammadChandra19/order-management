package product

import "database/sql"

//go:generate mockgen -source repository.go -destination mock/repository_mock.go -package=mock

// ProductSalesStats represents the statistics of product sales.
type ProductSalesStats struct {
	ProductName       string  `json:"product_name"`
	TotalQuantitySold int     `json:"total_quantity_sold"`
	TotalAmount       float64 `json:"total_amount"`
}

// ProductRepositoryInterface defines the interface for product-related database operations.
type ProductRepositoryInterface interface {
	// GetProductSalesStats retrieves the 5 most sold products along with their corresponding total quantities sold and total amounts.
	//
	// Returns:
	//
	// - ([]*ProductSalesStats): A slice of product sales statistics.
	//
	// - (error): An error, if any.
	GetProductSalesStats() ([]*ProductSalesStats, error)
}

type repository struct {
	db *sql.DB
}

func (r *repository) GetProductSalesStats() ([]*ProductSalesStats, error) {
	query := `
		SELECT 
			oi.product AS product_name,
			SUM(d.delivered_quantity) AS total_quantity_sold,
			SUM(oi.price_per_unit * d.delivered_quantity) AS total_amount
		FROM 
			order_items oi
		JOIN 
			deliveries d ON oi.id = d.order_item_id
		GROUP BY 
			oi.product
		ORDER BY 
			total_quantity_sold DESC, 
			total_amount DESC
		LIMIT 5
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productsSalesStats []*ProductSalesStats
	for rows.Next() {
		var productSalesStats ProductSalesStats
		err := rows.Scan(
			&productSalesStats.ProductName,
			&productSalesStats.TotalQuantitySold,
			&productSalesStats.TotalAmount,
		)
		if err != nil {
			return nil, err
		}
		productsSalesStats = append(productsSalesStats, &productSalesStats)
	}

	return productsSalesStats, nil
}

// NewProductRepository creates a new instance of ProductRepositoryInterface.
// Parameters:
// - db (*sql.DB): Database connection.
// Returns:
// - (ProductRepositoryInterface): New instance of ProductRepositoryInterface.
func NewProductRepository(db *sql.DB) ProductRepositoryInterface {
	return &repository{
		db: db,
	}
}
