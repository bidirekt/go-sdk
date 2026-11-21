package calls

import (
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("GET", "/test-only",
	bidirekt.Response(200, nil),
)
