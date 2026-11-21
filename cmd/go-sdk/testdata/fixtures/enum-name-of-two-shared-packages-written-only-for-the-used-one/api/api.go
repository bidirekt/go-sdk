package api

import (
	"example.com/shared/billing"
	"example.com/shared/orders"
	"github.com/bidirekt/go-sdk"
)

type Order struct {
	Status orders.Status `json:"status" bidirekt:"field"`
	Total  billing.Money `json:"total" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/orders/*",
	bidirekt.Response(200, Order{}),
)
