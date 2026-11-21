package api

import (
	"example.com/ticketing/internal/tickets"
	"github.com/gin-gonic/gin"
)

func NewRouter(service *tickets.Service) *gin.Engine {
	router := gin.Default()
	handler := &TicketHandler{service: service}
	group := router.Group("/api")
	group.POST("/tickets", handler.Book)
	group.GET("/tickets/:id", handler.Get)
	group.POST("/tickets/:id/cancel", handler.Cancel)
	group.GET("/tickets/stats", handler.Stats)
	return router
}
