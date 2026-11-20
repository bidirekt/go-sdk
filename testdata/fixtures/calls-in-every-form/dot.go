package petstore

import (
	. "github.com/bidirekt/go-sdk"
)

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusAdopted   PetStatus = "adopted"
)

var _ = Enum(PetStatusAvailable, PetStatusAdopted)

type PetColor string

var _ = Enum[PetColor]("black", "white")

type PetMood string

var _ = ((Enum[PetMood]))("calm", "playful")

var _ = ((Consumes))("petstore_reviews", "GET", "/ratings",
	((Response))(200, nil),
)
