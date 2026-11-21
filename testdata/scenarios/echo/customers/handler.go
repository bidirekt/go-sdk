package customers

import (
	"net/http"
	"sync"

	"example.com/shop/shared"
	"github.com/bidirekt/go-sdk"
	"github.com/labstack/echo/v4"
)

type Customer struct {
	ID    string `json:"id" bidirekt:"field"`
	Name  string `json:"name" bidirekt:"field"`
	Email string `json:"email" bidirekt:"field"`
	shared.Timestamps
}

type Handler struct {
	mutex     sync.Mutex
	customers map[string]Customer
}

func NewHandler() *Handler {
	return &Handler{customers: map[string]Customer{}}
}

var _ = bidirekt.Provides(http.MethodGet, "/customers/*",
	bidirekt.Response(http.StatusOK, Customer{}),
	bidirekt.Response(http.StatusNotFound, shared.Error{}),
)

func (this *Handler) Get(c echo.Context) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	customer, found := this.customers[c.Param("id")]
	if !found {
		return c.JSON(http.StatusNotFound, shared.Error{Message: "no such customer"})
	}
	return c.JSON(http.StatusOK, customer)
}
