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

var _ = bidirekt.Schema(PetStatus(""), bidirekt.Named("Status"))

var _ = bidirekt.Schema(Pet{}, bidirekt.Named("PetResponse"))

type Pet struct {
	Name   string    `json:"name" bidirekt:"field"`
	Status PetStatus `json:"status" bidirekt:"field"`
	Mother *Pet      `json:"mother" bidirekt:"field"`
}

type Owner struct {
	Pet  Pet    `json:"pet" bidirekt:"field"`
	Last string `json:"last" bidirekt:"field type=PetResponse"`
}

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(200, Pet{}),
)
