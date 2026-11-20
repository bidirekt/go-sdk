package compileerror

import (
	"github.com/bidirekt/go-sdk"
)

var count int = "many"

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(200, nil),
)
