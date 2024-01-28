package order

//go:generate mockgen -source service.go -destination mock/service_mock.go -package=mock

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MuhammadChandra19/order-management/internal/utils"
	"github.com/gin-gonic/gin"
)

// OrderServiceInterface defines the interface for order-related operations.
type OrderServiceInterface interface {
	// GetOrderList retrieves a list of orders based on the provided context.
	// Parameters:
	// - c (*gin.Context): Context for handling the request.
	GetOrderList(c *gin.Context)
}

type service struct {
	orderRepo OrderRepositoryInterface
}

func (s *service) GetOrderList(c *gin.Context) {

	// Extract parameters from the query string
	search := c.DefaultQuery("search", "")
	startDateStr := c.DefaultQuery("start_date", "2000-01-12T00:00:00Z")
	endDateStr := c.DefaultQuery("end_date", "2100-01-12T00:00:00Z")
	sortDirection := c.DefaultQuery("sort_direction", "DESC")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))   // Default limit to 5 if not provided
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0")) // Default offset to 0 if not provided

	startDate, err := utils.CompileDate(startDateStr)
	if err != nil {
		fmt.Println("err startDate", err)
	}
	endDate, err := utils.CompileDate(endDateStr)
	if err != nil {
		fmt.Println("err startDate", err)
	}

	orders, err := s.orderRepo.GetOrderList(search, startDate, endDate, sortDirection, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)

}

// NewOrderService creates a new OrderService instance with the provided OrderRepositoryInterface.
// Parameters:
// - orderRepo (OrderRepositoryInterface): Order repository interface implementation.
// Returns:
// - (OrderServiceInterface): New instance of OrderServiceInterface.
func NewOrderService(orderRepo OrderRepositoryInterface) OrderServiceInterface {
	return &service{
		orderRepo: orderRepo,
	}
}
