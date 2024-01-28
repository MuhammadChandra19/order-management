package order_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MuhammadChandra19/order-management/internal/order"
	"github.com/MuhammadChandra19/order-management/internal/order/mock"
	"github.com/MuhammadChandra19/order-management/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetOrderList(t *testing.T) {

	startDate, err := utils.CompileDate("2022-01-01T00:00:00Z")
	if err != nil {
		fmt.Println("err startDate", err)
	}
	endDate, err := utils.CompileDate("2022-12-31T23:59:59Z")
	if err != nil {
		fmt.Println("err startDate", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Create a new instance of your service
	orderRepo := mock.NewMockOrderRepositoryInterface(ctrl)
	service := order.NewOrderService(orderRepo)

	// Create a new Gin router
	router := gin.New()
	router.GET("/api/orders", service.GetOrderList)

	// Define test cases
	testCases := []struct {
		name        string
		queryParams string
		mockFn      func(*mock.MockOrderRepositoryInterface)

		expectedStatus int
	}{
		{
			name:        "Successful request",
			queryParams: "?search=test&start_date=2022-01-01T00:00:00Z&end_date=2022-12-31T23:59:59Z&sort_direction=DESC&limit=5&offset=0",
			mockFn: func(mori *mock.MockOrderRepositoryInterface) {
				mori.EXPECT().GetOrderList("test", startDate, endDate, "DESC", 5, 0).Return([]*order.OrderInfo{
					{
						Id:           "1",
						OrderName:    "OrderName",
						CompanyName:  "CompanyName",
						CustomerName: "CustomerName",
						ProductName:  "ProductName",
						OrderDate:    startDate,
						Delivered:    "5",
						TotalAmount:  11.5,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Error request",
			queryParams: "?search=test&start_date=2022-01-01T00:00:00Z&end_date=2022-12-31T23:59:59Z&sort_direction=DESC&limit=5&offset=0",
			mockFn: func(mori *mock.MockOrderRepositoryInterface) {
				mori.EXPECT().
					GetOrderList("test", startDate, endDate, "DESC", 5, 0).
					Return(nil, errors.New("unexpected error when getting orders"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new HTTP request with the specified query parameters
			req, err := http.NewRequest(http.MethodGet, "/api/orders"+tc.queryParams, nil)
			tc.mockFn(orderRepo)
			assert.NoError(t, err)

			// Create a response recorder to record the response
			res := httptest.NewRecorder()

			// Serve the HTTP request using the Gin router
			router.ServeHTTP(res, req)

			// Validate the response status code
			assert.Equal(t, tc.expectedStatus, res.Code)

		})
	}
}
