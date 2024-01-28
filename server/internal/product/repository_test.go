package product_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MuhammadChandra19/order-management/internal/product"
	"github.com/stretchr/testify/assert"
)

func GetRowsMock(sm sqlmock.Sqlmock) *sqlmock.Rows {
	return sm.NewRows([]string{"product_name", "total_quantity_sold", "total_amount"})
}

func TestProduct_GetProductSalesStats(t *testing.T) {
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

	testCases := []struct {
		name       string
		mockFunc   func(sqlmock.Sqlmock)
		assertFunc func(*testing.T, []*product.ProductSalesStats, error)
	}{
		{
			name: "query error",
			mockFunc: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(sql.ErrConnDone)
			},
			assertFunc: func(t *testing.T, pss []*product.ProductSalesStats, err error) {
				assert.Equal(t, sql.ErrConnDone.Error(), err.Error())
			},
		},
		{
			name: "scan error",
			mockFunc: func(s sqlmock.Sqlmock) {
				rows := GetRowsMock(s)
				rows.AddRow("Hand sanitizer", "delivered", "amount")
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(rows)
			},
			assertFunc: func(t *testing.T, pss []*product.ProductSalesStats, err error) {
				assert.Contains(t, err.Error(), "Scan error")
			},
		},
		{
			name: "success",
			mockFunc: func(s sqlmock.Sqlmock) {
				rows := GetRowsMock(s)
				rows.AddRow("Hand sanitizer", 135, 12295.8598)
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(rows)
			},
			assertFunc: func(t *testing.T, pss []*product.ProductSalesStats, err error) {
				assert.Equal(t, &product.ProductSalesStats{
					ProductName:       "Hand sanitizer",
					TotalQuantitySold: 135,
					TotalAmount:       12295.8598,
				}, pss[0])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			defer db.Close()

			tc.mockFunc(mock)
			productRepo := product.NewProductRepository(db)
			list, err := productRepo.GetProductSalesStats()
			fmt.Println("here", err)
			tc.assertFunc(t, list, err)
		})
	}
}
