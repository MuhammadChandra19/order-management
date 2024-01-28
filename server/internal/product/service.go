package product

//go:generate mockgen -source service.go -destination mock/service_mock.go -package=mock
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProductServiceInterface defines the interface for product-related services.
type ProductServiceInterface interface {
	// GetProductSalesStats retrieves product sales statistics.
	//
	// Parameters:
	//
	// - c (*gin.Context): Context object for Gin HTTP request.
	//
	// Returns:
	// None.
	GetProductSalesStats(c *gin.Context)
}

type service struct {
	productRepo ProductRepositoryInterface
}

func (s *service) GetProductSalesStats(c *gin.Context) {
	products, err := s.productRepo.GetProductSalesStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// NewProductService creates a new instance of ProductServiceInterface.
//
// Parameters:
//
// - productRepo (ProductRepositoryInterface): Product repository interface.
//
// Returns:
//
// - (ProductServiceInterface): New instance of ProductServiceInterface.
func NewProductService(productRepo ProductRepositoryInterface) ProductServiceInterface {
	return &service{
		productRepo: productRepo,
	}
}
