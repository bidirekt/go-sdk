package petstore

import (
	"net/http"

	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides(http.MethodDelete, "/pets/*",
	bidirekt.Response(http.StatusNoContent, nil),
)
