package petstore

import (
	"github.com/bidirekt/go-sdk"
)

func refund() {
	_ = bidirekt.Consumes("payments_api", "GET", "/charges/*",
		bidirekt.Response(410, nil),
	)
}
