package order_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestOrder_GetOrderList(t *testing.T) {

	testCases := map[string]struct {
		mockFunc   func(sqlmock.Sqlmock)
		assertFunc func(*testing.T, error)
	}{
		"success": {},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tc.mockFunc(mock)
			// orderRepo := order.NewOrderRepository(db)
			// err = orderRepo.SeedData()
			// fmt.Println(err)
			tc.assertFunc(t, err)
		})
	}
}
