//go:build never

package calls

import (
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("GET", "/never",
	bidirekt.Response(200, nil),
)
