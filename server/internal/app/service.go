package app

import (
	"github.com/MuhammadChandra19/order-management/internal/db"
	"github.com/MuhammadChandra19/order-management/internal/order"
)

type App struct {
	OrderService order.OrderServiceInterface
}

func InitApp() *App {
	pg := db.NewDatabase()
	// _ = customer.NewCustomerRepository(pg)
	orderRepo := order.NewOrderRepository(&pg)
	// deliveryRepo := delivery.NewDeliveryRepository(pg)
	orderService := order.NewOrderService(orderRepo)

	return &App{
		OrderService: orderService,
	}
}
