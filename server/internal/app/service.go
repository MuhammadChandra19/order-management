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

// InitApp initializes the application by creating instances of database connections,
// repositories, and services for orders and products. It sets up the necessary
// dependencies and returns an instance of the App struct containing the initialized
// services for handling orders and products.
//
// Parameters:
// - None
//
// Returns:
// - *App: An instance of the App struct containing initialized order and product services.
//
// Notes:
// - The function creates a new database connection using db.NewDatabase().
// - It initializes repositories for orders and products.
// - It creates services for handling orders and products.
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
