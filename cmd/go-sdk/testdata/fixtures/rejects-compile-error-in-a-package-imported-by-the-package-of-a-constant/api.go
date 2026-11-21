package app

import (
	"net/http"

	"example.com/app/routes"
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides(http.MethodGet, routes.Pets,
	bidirekt.Response(http.StatusOK, nil),
)

var _ = bidirekt.Provides(http.MethodGet, "/health",
	bidirekt.Response(http.StatusOK, nil),
)
