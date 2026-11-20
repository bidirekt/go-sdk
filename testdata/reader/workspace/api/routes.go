package api

import (
	"example.com/shared"
	"github.com/bidirekt/go-sdk"
)

var status = shared.PetStatusAvailable

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(200, nil),
)
