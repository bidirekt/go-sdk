package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

var _ = bidirekt.Schema(Pet{}, bidirekt.Named("PetResponse"))

var _ = bidirekt.Schema(Pet{}, bidirekt.Named("PetDetail"))
