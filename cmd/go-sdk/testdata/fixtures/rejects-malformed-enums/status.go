package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusAdopted   PetStatus = "adopted"
)

type PetSize string

type PetColor string

type PetMood string

var mood = PetMood("calm")

var _ = bidirekt.Enum(PetStatusAvailable, PetStatusAdopted, PetStatusAvailable)

var _ = bidirekt.Enum(PetSize("small"), PetSize(""))

var _ = bidirekt.Enum[PetColor]()

var _ = bidirekt.Enum(mood)

var _ = bidirekt.Enum("small", "large")

type PetBreed string

var breeds = []PetBreed{"poodle", "beagle"}

var _ = bidirekt.Enum(breeds...)
