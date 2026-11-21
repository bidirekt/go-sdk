package app

import (
	"net/http"

	"example.com/money"
	"github.com/bidirekt/go-sdk"
)

var _ money.Money

type Order struct {
	Totals []any `json:"totals" bidirekt:"field type=array items=Money"`
}

var _ = bidirekt.Provides(http.MethodGet, "/orders/*",
	bidirekt.Response(http.StatusOK, Order{}),
)
