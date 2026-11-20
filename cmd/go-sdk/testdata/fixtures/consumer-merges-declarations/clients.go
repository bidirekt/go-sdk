package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Charge struct {
	ID string `json:"id" bidirekt:"field"`
}

type ChargeSummary struct {
	Total int `json:"total" bidirekt:"field"`
}

type NewCharge struct {
	Amount int `json:"amount" bidirekt:"field"`
}

type RetriedCharge struct {
	Attempt int `json:"attempt" bidirekt:"field"`
}

var _ = bidirekt.Consumes("payments_api", "GET", "/charges/*",
	bidirekt.Response(200, Charge{}),
)

var _ = bidirekt.Consumes("payments_api", "GET", "/charges/*",
	bidirekt.Response(200, Charge{}),
	bidirekt.Response(404, nil),
)
