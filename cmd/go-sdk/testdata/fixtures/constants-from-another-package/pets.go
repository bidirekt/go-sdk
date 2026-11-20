package petstore

import (
	"net/http"

	"example.com/petstore/routes"
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides(http.MethodGet, routes.Pets,
	bidirekt.Response(http.StatusOK, nil),
	bidirekt.Response(routes.Gone, nil),
)

var _ = bidirekt.Consumes(routes.ReviewsAPI, http.MethodGet, routes.Pets+"/*/ratings",
	bidirekt.Response(http.StatusOK, nil),
)
