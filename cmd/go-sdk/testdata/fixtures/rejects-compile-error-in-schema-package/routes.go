package petstore

import (
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("GET", "/health",
	bidirekt.Response(204, nil),
)
