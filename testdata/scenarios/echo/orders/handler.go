package orders

import (
	"fmt"
	"net/http"
	"sync"

	"example.com/shop/shared"
	"github.com/bidirekt/go-sdk"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	mutex  sync.Mutex
	orders map[string]Order
}

func NewHandler() *Handler {
	return &Handler{orders: map[string]Order{}}
}

func (this *Handler) Register(group *echo.Group) {
	group.POST("", this.Place)
	group.GET("/:id", this.Get)
	group.PUT("/:id/status", this.ChangeStatus)
}

var _ = bidirekt.Provides(http.MethodPost, "/orders",
	bidirekt.Request(PlaceOrder{}),
	bidirekt.Response(http.StatusCreated, Order{}),
	bidirekt.Response(http.StatusBadRequest, shared.Error{}),
)

func (this *Handler) Place(c echo.Context) error {
	var request struct {
		CustomerID string `json:"customerId"`
		Lines      []Line `json:"lines"`
	}
	if err := c.Bind(&request); err != nil || len(request.Lines) == 0 {
		return c.JSON(http.StatusBadRequest, shared.Error{Message: "an order needs lines"})
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	order := Order{
		ID:         fmt.Sprintf("ord_%d", len(this.orders)+1),
		CustomerID: request.CustomerID,
		Status:     StatusPlaced,
		Lines:      request.Lines,
		Payment:    CardPayment{Brand: "visa", LastFour: "4242"},
		Timestamps: shared.NewTimestamps(),
	}
	this.orders[order.ID] = order
	return c.JSON(http.StatusCreated, order)
}

var _ = bidirekt.Provides(http.MethodGet, "/orders/*",
	bidirekt.Response(http.StatusOK, Order{}),
	bidirekt.Response(http.StatusNotFound, shared.Error{}),
)

func (this *Handler) Get(c echo.Context) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	order, found := this.orders[c.Param("id")]
	if !found {
		return c.JSON(http.StatusNotFound, shared.Error{Message: "no such order"})
	}
	return c.JSON(http.StatusOK, order)
}

type StatusChange struct {
	Status Status `json:"status" bidirekt:"field"`
}

var _ = bidirekt.Provides(http.MethodPut, "/orders/*/status",
	bidirekt.Request(StatusChange{}),
	bidirekt.Response(http.StatusOK, Order{}),
	bidirekt.Response(http.StatusNotFound, shared.Error{}),
)

func (this *Handler) ChangeStatus(c echo.Context) error {
	var change StatusChange
	if err := c.Bind(&change); err != nil {
		return c.JSON(http.StatusBadRequest, shared.Error{Message: err.Error()})
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	order, found := this.orders[c.Param("id")]
	if !found {
		return c.JSON(http.StatusNotFound, shared.Error{Message: "no such order"})
	}
	order.Status = change.Status
	order.Touch()
	this.orders[order.ID] = order
	return c.JSON(http.StatusOK, order)
}
