package petstore

import (
	"example.com/shared"
	"github.com/bidirekt/go-sdk"
)

type LineItem struct {
	Sku      string `json:"sku" bidirekt:"field"`
	Quantity int    `json:"quantity" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/cart",
	bidirekt.Response(200, []LineItem{}),
)

var _ = bidirekt.Provides("GET", "/customer",
	bidirekt.Response(200, shared.Customer{}),
)
