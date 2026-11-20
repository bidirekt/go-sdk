package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Charge struct {
	ID string `json:"id" bidirekt:"field"`
}

type NewCharge struct {
	Amount float64 `json:"amount" bidirekt:"field"`
}

func charge() {
	_ = bidirekt.Consumes("payments_api", "POST", "/charges",
		bidirekt.Request(NewCharge{}),
		bidirekt.Response(201, Charge{}),
		bidirekt.Response(402, nil),
	)
}

func retry() {
	_ = bidirekt.Consumes("payments_api", "POST", "/charges",
		bidirekt.Request(NewCharge{}),
		bidirekt.Response(201, Charge{}),
		bidirekt.Response(409, nil),
	)
}
