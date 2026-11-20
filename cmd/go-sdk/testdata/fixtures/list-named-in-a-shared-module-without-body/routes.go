package petstore

import (
	"example.com/shared"
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(200, shared.Pet{}),
)
