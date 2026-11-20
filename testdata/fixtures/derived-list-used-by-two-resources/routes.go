package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(200, []Pet{}),
)

var _ = bidirekt.Provides("GET", "/pets/adopted",
	bidirekt.Response(200, []Pet{}),
)
