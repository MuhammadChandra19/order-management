package product_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MuhammadChandra19/order-management/internal/product"
	"github.com/MuhammadChandra19/order-management/internal/product/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetProductSalesStats(t *testing.T) {

	ctrl := gomock.NewController(t)

	productRepo := mock.NewMockProductRepositoryInterface(ctrl)
	service := product.NewProductService(productRepo)

	// Create a new Gin router
	router := gin.New()
	router.GET("/api/product-sale-stats", service.GetProductSalesStats)

	testCases := []struct {
		name   string
		mockFn func(*mock.MockProductRepositoryInterface)

		expectedStatus int
	}{
		{
			name: "success request",
			mockFn: func(mpri *mock.MockProductRepositoryInterface) {
				mpri.EXPECT().GetProductSalesStats().Return([]*product.ProductSalesStats{
					{
						ProductName:       "Hand sanitizer",
						TotalQuantitySold: 135,
						TotalAmount:       12295.8598,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error request",
			mockFn: func(mpri *mock.MockProductRepositoryInterface) {
				mpri.EXPECT().GetProductSalesStats().Return(nil, errors.New("unexpected error when getting stats"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new HTTP request with the specified query parameters
			req, err := http.NewRequest(http.MethodGet, "/api/product-sale-stats", nil)
			tc.mockFn(productRepo)
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
