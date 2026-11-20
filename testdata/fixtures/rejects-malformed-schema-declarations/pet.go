package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

type Owner struct {
	Name string `json:"name"`
}

type Breed string

var _ = bidirekt.Schema(Pet{})

var _ = bidirekt.Schema(Owner{}, bidirekt.Named("Keeper"))

var _ = bidirekt.Schema(Breed(""), bidirekt.Named("PetBreed"))

var _ = bidirekt.Schema("", bidirekt.Named("Text"))
