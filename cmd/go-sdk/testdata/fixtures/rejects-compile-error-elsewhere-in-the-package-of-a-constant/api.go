package app

import (
	"net/http"

	"example.com/app/routes"
	"github.com/bidirekt/go-sdk"
)

const healthPath = routes.Health

var _ = bidirekt.Provides(http.MethodGet, healthPath,
	bidirekt.Response(http.StatusOK, nil),
)
