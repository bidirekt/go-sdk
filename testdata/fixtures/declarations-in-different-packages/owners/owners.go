package owners

import (
	"example.com/petstore/pets"
	"github.com/bidirekt/go-sdk"
)

type Owner struct {
	Name string     `json:"name" bidirekt:"field"`
	Pets []pets.Pet `json:"pets" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/owners/*",
	bidirekt.Response(200, Owner{}),
)
