package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusOnHold    PetStatus = "on_hold"
)

var _ = bidirekt.Enum(PetStatusAvailable, PetStatusOnHold)

type Pet struct {
	Status         PetStatus  `json:"status" bidirekt:"field"`
	PreviousStatus *PetStatus `json:"previousStatus" bidirekt:"field"`
}
