package petstore

import (
	"net/http"

	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Consumes("petstore_api", http.MethodGet, "/pets",
	bidirekt.Response(http.StatusOK, []Pet{}),
)

var _ = bidirekt.Consumes("petstore_api", http.MethodPost, "/pets",
	bidirekt.Request(NewPet{}),
	bidirekt.Response(http.StatusCreated, Pet{}),
	bidirekt.Response(http.StatusBadRequest, nil),
)

var _ = bidirekt.Consumes("reviews_api", http.MethodGet, "/ratings/*",
	bidirekt.Response(http.StatusOK, map[string]float64{}, bidirekt.Named("Ratings")),
)
