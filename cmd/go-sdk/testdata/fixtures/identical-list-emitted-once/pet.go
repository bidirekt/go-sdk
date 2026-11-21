package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/pets/featured",
	bidirekt.Response(200, []Pet{}, bidirekt.Named("PetCollection")),
)

var _ = bidirekt.Consumes("shelter_api", "GET", "/residents",
	bidirekt.Response(200, []Pet{}, bidirekt.Named("PetCollection")),
)
