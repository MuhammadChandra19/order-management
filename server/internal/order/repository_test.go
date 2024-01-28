package order_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MuhammadChandra19/order-management/internal/order"
	"github.com/stretchr/testify/assert"
)

func GetRowsMock(sm sqlmock.Sqlmock) *sqlmock.Rows {
	return sm.NewRows([]string{"id", "order_name", "company_name", "customer_name", "product_name", "order_date", "delivered", "total_amount"})
}

func TestOrder_GetOrderList(t *testing.T) {
	now := time.Now()
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

	testCases := map[string]struct {
		search        string
		startDate     time.Time
		endDate       time.Time
		sortDirection string
		limit         int
		offset        int

		mockFunc   func(sqlmock.Sqlmock, string, time.Time, time.Time, int, int)
		assertFunc func(*testing.T, []*order.OrderInfo, error)
	}{
		"query error": {
			search:        "",
			startDate:     now,
			endDate:       now,
			sortDirection: "DESC",
			limit:         5,
			offset:        0,

			mockFunc: func(s sqlmock.Sqlmock, search string, startDate, endDate time.Time, limit, offset int) {
				newQuery := query
				newQuery += `
					AND 
						o.created_at BETWEEN $1::timestamp AND $2::timestamp
					GROUP BY
						o.id, cc.company_id, c.id, oi.product
					ORDER BY
						o.created_at DESC
					LIMIT $3 OFFSET $4
				`
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(startDate, endDate, limit, offset).
					WillReturnError(sql.ErrConnDone)
			},
			assertFunc: func(t *testing.T, oi []*order.OrderInfo, err error) {
				assert.Equal(t, sql.ErrConnDone.Error(), err.Error())
			},
		},
		"scan error": {
			search:        "",
			startDate:     now,
			endDate:       now,
			sortDirection: "DESC",
			limit:         5,
			offset:        0,

			mockFunc: func(s sqlmock.Sqlmock, search string, startDate, endDate time.Time, limit, offset int) {
				newQuery := query
				newQuery += `
					AND 
						o.created_at BETWEEN $1::timestamp AND $2::timestamp
					GROUP BY
						o.id, cc.company_id, c.id, oi.product
					ORDER BY
						o.created_at DESC
					LIMIT $3 OFFSET $4
				`
				rows := GetRowsMock(s)
				rows.AddRow("1", "OrderName", "CompanyName", "CustomerName", "ProductName", "OrderDate", "Delivered", "TotalAmount")
				s.ExpectQuery(regexp.QuoteMeta(newQuery)).
					WithArgs(startDate, endDate, limit, offset).
					WillReturnRows(rows)
			},
			assertFunc: func(t *testing.T, oi []*order.OrderInfo, err error) {
				assert.Contains(t, err.Error(), "Scan error")
			},
		},
		"success": {
			search:        "",
			startDate:     now,
			endDate:       now,
			sortDirection: "DESC",
			limit:         5,
			offset:        0,

			mockFunc: func(s sqlmock.Sqlmock, search string, startDate, endDate time.Time, limit, offset int) {
				newQuery := query
				newQuery += `
					AND 
						o.created_at BETWEEN $1::timestamp AND $2::timestamp
					GROUP BY
						o.id, cc.company_id, c.id, oi.product
					ORDER BY
						o.created_at DESC
					LIMIT $3 OFFSET $4
				`
				rows := GetRowsMock(s)
				rows.AddRow("1", "OrderName", "CompanyName", "CustomerName", "ProductName", now, "5", 11.5)
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(startDate, endDate, limit, offset).
					WillReturnRows(rows)
			},
			assertFunc: func(t *testing.T, oi []*order.OrderInfo, err error) {
				assert.Equal(t, &order.OrderInfo{
					Id:           "1",
					OrderName:    "OrderName",
					CompanyName:  "CompanyName",
					CustomerName: "CustomerName",
					ProductName:  "ProductName",
					OrderDate:    now,
					Delivered:    "5",
					TotalAmount:  11.5,
				}, oi[0])
			},
		},
		"should use search params and ASC order": {
			search:        "Box",
			startDate:     now,
			endDate:       now,
			sortDirection: "ASC",
			limit:         5,
			offset:        0,

			mockFunc: func(s sqlmock.Sqlmock, search string, startDate, endDate time.Time, limit, offset int) {
				newQuery := query
				newQuery += `
					AND 
						(o.order_name ILIKE $1 OR oi.product ILIKE $1) 
					AND 
						o.created_at BETWEEN $2::timestamp AND $3::timestamp
					GROUP BY
						o.id, cc.company_id, c.id, oi.product
					ORDER BY
						o.created_at ASC
					LIMIT $4 OFFSET $5
				`
				rows := GetRowsMock(s)
				rows.AddRow("1", "OrderName", "CompanyName", "CustomerName", "ProductName", now, "5", 11.5)
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("%"+search+"%", startDate, endDate, limit, offset).
					WillReturnRows(rows)
			},
			assertFunc: func(t *testing.T, oi []*order.OrderInfo, err error) {
				assert.Equal(t, &order.OrderInfo{
					Id:           "1",
					OrderName:    "OrderName",
					CompanyName:  "CompanyName",
					CustomerName: "CustomerName",
					ProductName:  "ProductName",
					OrderDate:    now,
					Delivered:    "5",
					TotalAmount:  11.5,
				}, oi[0])
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			defer db.Close()

			tc.mockFunc(mock, tc.search, tc.startDate, tc.endDate, tc.limit, tc.offset)
			orderRepo := order.NewOrderRepository(db)
			list, err := orderRepo.GetOrderList(tc.search, tc.startDate, tc.endDate, tc.sortDirection, tc.limit, tc.offset)
			fmt.Println(err)
			tc.assertFunc(t, list, err)
		})
	}
}
