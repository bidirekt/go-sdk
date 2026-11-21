package app

import (
	"net/http"

	"example.com/app/routes"
	"github.com/bidirekt/go-sdk"
)

const pets = routes.Base + "/pets"

var _ = bidirekt.Provides(http.MethodGet, pets,
	bidirekt.Response(http.StatusOK, nil),
)

var _ = bidirekt.Provides(http.MethodGet, routes.Base+"/owners",
	bidirekt.Response(http.StatusOK, nil),
)
