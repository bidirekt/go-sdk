package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusAdopted   PetStatus = "adopted"
)

var _ = bidirekt.Enum(PetStatusAvailable, PetStatusAdopted)

type Pet struct {
	Status PetStatus `json:"status" bidirekt:"field"`
}
