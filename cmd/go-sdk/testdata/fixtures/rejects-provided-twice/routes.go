package petstore

import (
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(200, nil),
)

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(404, nil),
)
