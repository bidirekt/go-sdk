package app

import (
	"net/http"

	"example.com/app/kinds"
	"github.com/bidirekt/go-sdk"
)

type Ticket struct {
	Tags []kinds.Tag `json:"tags" bidirekt:"field"`
}

var _ = bidirekt.Provides(http.MethodGet, "/tickets/*",
	bidirekt.Response(http.StatusOK, Ticket{}),
)

var _ = bidirekt.Provides(http.MethodGet, "/tags",
	bidirekt.Response(http.StatusOK, kinds.Tags{}),
)
