package api

import (
	"example.com/shared"
	"github.com/bidirekt/go-sdk"
)

type Pet struct {
	Status shared.PetStatus `json:"status" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(200, Pet{}),
)
