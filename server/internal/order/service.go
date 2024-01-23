package order

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderServiceInterface interface {
	GetOrderList(c *gin.Context)
}

type service struct {
	orderRepo OrderRepositoryInterface
}

func (s *service) GetOrderList(c *gin.Context) {

	// Extract parameters from the query string
	search := c.DefaultQuery("search", "")
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))   // Default limit to 5 if not provided
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0")) // Default offset to 0 if not provided

	// Parse startDate and endDate as time.Time
	startDate, _ := time.Parse(time.RFC3339, startDateStr)
	endDate, _ := time.Parse(time.RFC3339, endDateStr)

	orders, err := s.orderRepo.GetOrderList(search, startDate, endDate, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)

}

func NewOrderService(orderRepo OrderRepositoryInterface) OrderServiceInterface {
	return &service{
		orderRepo: orderRepo,
	}
}
