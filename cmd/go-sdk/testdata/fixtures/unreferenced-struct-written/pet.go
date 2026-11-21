package petstore

import (
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("GET", "/health",
	bidirekt.Response(204, nil),
)

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

type Draft struct {
	Name string `json:"name"`
}
