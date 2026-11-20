package nonconstant

import (
	"github.com/bidirekt/go-sdk"
)

var petsPath = "/pets"

var _ = bidirekt.Provides("GET", petsPath,
	bidirekt.Response(200, nil),
)
