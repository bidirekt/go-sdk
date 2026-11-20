package petstore

import (
	"net/http"

	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("DELETE", "/pets/*",
	bidirekt.Response(http.StatusNoContent, nil),
	bidirekt.Response(http.StatusNotFound, nil),
)

var _ = bidirekt.Provides(http.MethodGet, "/health",
	bidirekt.Response(200, nil),
)
