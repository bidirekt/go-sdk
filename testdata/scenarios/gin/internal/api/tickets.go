package api

import (
	"errors"
	"net/http"
	"strconv"

	"example.com/ticketing/internal/tickets"
	"github.com/bidirekt/go-sdk"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	service *tickets.Service
}

var _ = bidirekt.Provides("POST", "/api/tickets",
	bidirekt.Request(BookTicketRequest{}),
	bidirekt.Response(http.StatusCreated, TicketResponse{}),
	bidirekt.Response(http.StatusBadRequest, ErrorResponse{}),
	bidirekt.Response(http.StatusConflict, ErrorResponse{}),
)

func (this *TicketHandler) Book(c *gin.Context) {
	var request BookTicketRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	ticket, err := this.service.Book(request.EventID, request.Quantity, request.Email, request.Coupon)
	if errors.Is(err, tickets.ErrSoldOut) {
		c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ticketToResponse(ticket))
}

var _ = bidirekt.Provides("GET", "/api/tickets/*",
	bidirekt.Response(http.StatusOK, TicketResponse{}),
	bidirekt.Response(http.StatusNotFound, ErrorResponse{}),
)

func (this *TicketHandler) Get(c *gin.Context) {
	ticket, err := this.find(c)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ticketToResponse(ticket))
}

var _ = bidirekt.Provides("POST", "/api/tickets/*/cancel",
	bidirekt.Request(nil),
	bidirekt.Response(http.StatusOK, TicketResponse{}),
	bidirekt.Response(http.StatusNotFound, ErrorResponse{}),
)

func (this *TicketHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: tickets.ErrNotFound.Error()})
		return
	}
	ticket, err := this.service.Cancel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ticketToResponse(ticket))
}

var _ = bidirekt.Provides("GET", "/api/tickets/stats",
	bidirekt.Response(http.StatusOK, struct {
		Active   int `json:"active" bidirekt:"field"`
		Canceled int `json:"canceled" bidirekt:"field"`
	}{}, bidirekt.Named("TicketStats")),
)

func (this *TicketHandler) Stats(c *gin.Context) {
	active, canceled := this.service.Counts()
	c.JSON(http.StatusOK, gin.H{"active": active, "canceled": canceled})
}

func (this *TicketHandler) find(c *gin.Context) (tickets.Ticket, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return tickets.Ticket{}, tickets.ErrNotFound
	}
	return this.service.Find(id)
}
