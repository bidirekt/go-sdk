package petstore

import (
	"net/http"

	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(http.StatusOK, Pet{}),
	bidirekt.Response(200, nil),
)
