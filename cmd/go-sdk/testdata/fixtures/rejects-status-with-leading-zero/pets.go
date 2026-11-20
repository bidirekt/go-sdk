package petstore

import (
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(0200, nil),
)
