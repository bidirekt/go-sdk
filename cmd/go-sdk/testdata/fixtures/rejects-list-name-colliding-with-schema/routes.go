package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

type PetList struct {
	Total int `json:"total" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(200, []Pet{}),
)

var _ = bidirekt.Provides("GET", "/pets/featured",
	bidirekt.Response(200, []Pet{}, bidirekt.Named("Pet")),
)
