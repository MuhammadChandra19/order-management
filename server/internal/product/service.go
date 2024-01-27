package product

//go:generate mockgen -source service.go -destination mock/service_mock.go -package=mock
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductServiceInterface interface {
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

func NewProductService(productRepo ProductRepositoryInterface) ProductServiceInterface {
	return &service{
		productRepo: productRepo,
	}
}
