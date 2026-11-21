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

type Receipt struct {
	ID string `json:"id" bidirekt:"field"`
}

type AuditEntry struct {
	Message string `json:"message" bidirekt:"field"`
}

var _ = bidirekt.Consumes("payments_api", "POST", "/charges",
	bidirekt.Request(NewCharge{}),
	bidirekt.Response(201, Charge{}),
	bidirekt.Response(402, nil),
)

var _ = bidirekt.Consumes("payments_api", "GET", "/charges",
	bidirekt.Response(200, []Charge{}),
)

var _ = bidirekt.Consumes("payments_api", "GET", "/charges/recent",
	bidirekt.Response(200, []Charge{}, bidirekt.Named("ChargeCollection")),
)

var _ = bidirekt.Consumes("payments_api", "GET", "/charges/*/receipt",
	bidirekt.Response(200, Receipt{}),
)

var _ = bidirekt.Consumes("payments_api", "PATCH", "/charges/*",
	bidirekt.Request(nil),
	bidirekt.Response(200, Charge{}),
	bidirekt.Response(404, nil),
)

var _ = bidirekt.Consumes("audit_api", "PUT", "/entries/*",
	bidirekt.Request(AuditEntry{}),
	bidirekt.Response(204, nil),
)
