package main

import (
	"example.com/shop/customers"
	"example.com/shop/orders"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	server := echo.New()
	server.Use(middleware.Logger(), middleware.Recover())
	orders.NewHandler().Register(server.Group("/orders"))
	server.GET("/customers/:id", customers.NewHandler().Get)
	server.Logger.Fatal(server.Start(":8080"))
}
