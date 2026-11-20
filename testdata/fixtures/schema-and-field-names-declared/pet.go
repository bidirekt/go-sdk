package petstore

import (
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Schema(Pet{}, bidirekt.Named("PetResponse"))

type Pet struct {
	BirthDate string `json:"birthDate" bidirekt:"field name=birth_date"`
}
