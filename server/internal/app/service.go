package app

import (
	"github.com/MuhammadChandra19/order-management/internal/db"
	"github.com/MuhammadChandra19/order-management/internal/order"
	"github.com/MuhammadChandra19/order-management/internal/product"
)

type App struct {
	OrderService   order.OrderServiceInterface
	ProductService product.ProductServiceInterface
}

func InitApp() *App {
	pg := db.NewDatabase()
	// _ = customer.NewCustomerRepository(pg)
	orderRepo := order.NewOrderRepository(pg)
	// deliveryRepo := delivery.NewDeliveryRepository(pg)
	orderService := order.NewOrderService(orderRepo)

	productRepo := product.NewProductRepository(pg)
	productService := product.NewProductService(productRepo)

	return &App{
		OrderService:   orderService,
		ProductService: productService,
	}
}
