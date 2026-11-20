package routes

import (
	"net/http"

	"github.com/bidirekt/go-sdk"
)

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

var _ = bidirekt.Provides(http.MethodGet, "/pets/*",
	bidirekt.Response(http.StatusOK, Pet{}),
)
