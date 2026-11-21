package app

import (
	"net/http"

	"example.com/app/kinds"
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Enum(kinds.Open, kinds.Closed)

type Ticket struct {
	Status kinds.Status `json:"status" bidirekt:"field"`
}

var _ = bidirekt.Provides(http.MethodGet, "/tickets/*",
	bidirekt.Response(http.StatusOK, Ticket{}),
)
