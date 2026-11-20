package petstore

import (
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Schema(Owner{}, bidirekt.Named("PetOwner"))

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pet struct {
	Owner          Owner   `json:"owner" bidirekt:"field"`
	PreviousOwners []Owner `json:"previousOwners" bidirekt:"field"`
}
