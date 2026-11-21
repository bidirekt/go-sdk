package shared

import (
	"github.com/bidirekt/go-sdk"
)

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pets []Pet

var _ = bidirekt.Schema(Pets{}, bidirekt.Named("PetCollection"))
