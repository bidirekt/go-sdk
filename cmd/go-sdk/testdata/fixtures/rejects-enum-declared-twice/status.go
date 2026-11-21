package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusAdopted   PetStatus = "adopted"
)

var _ = bidirekt.Enum(PetStatusAvailable)

var _ = bidirekt.Enum(PetStatusAvailable, PetStatusAdopted)
