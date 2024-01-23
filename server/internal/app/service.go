package app

import (
	"github.com/MuhammadChandra19/order-management/internal/db"

	"github.com/MuhammadChandra19/order-management/internal/customer"
	"github.com/MuhammadChandra19/order-management/internal/delivery"
	"github.com/MuhammadChandra19/order-management/internal/order"
)

type App struct {
	OrderService order.OrderServiceInterface
}

func InitApp() *App {
	pg, isPopulated := db.NewDatabase()
	customerRepo := customer.NewCustomerRepository(pg)
	orderRepo := order.NewOrderRepository(pg)
	deliveryRepo := delivery.NewDeliveryRepository(pg)

	if !isPopulated {
		customerRepo.SeedData()
		orderRepo.SeedData()
		deliveryRepo.SeedData()
	}

	orderService := order.NewOrderService(orderRepo)

	return &App{
		OrderService: orderService,
	}
}
