package main

import (
	"time"

	// "github.com/MuhammadChandra19/order-management/internal/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/MuhammadChandra19/order-management/internal/app"
)

func main() {
	app := app.InitApp()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api/orders", app.OrderService.GetOrderList)
	r.GET("/api/product-sale-stats", app.ProductService.GetProductSalesStats)
	r.Run((":8080"))
}
