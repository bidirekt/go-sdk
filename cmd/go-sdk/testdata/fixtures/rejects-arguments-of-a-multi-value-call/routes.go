package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

type PetStatus string

func route() (string, string) {
	return "GET", "/pets"
}

func consumed() (string, string, string) {
	return "petstore_reviews", "GET", "/ratings"
}

func ok() (int, any) {
	return 200, Pet{}
}

func body() (any, bidirekt.Name) {
	return Pet{}, bidirekt.Name{}
}

func statuses() (PetStatus, PetStatus) {
	return "available", "adopted"
}

var _ = bidirekt.Provides(route())

var _ = bidirekt.Consumes(consumed())

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(ok()),
)

var _ = bidirekt.Provides("POST", "/pets",
	bidirekt.Request(body()),
	bidirekt.Response(201, nil),
)

var _ = bidirekt.Schema(body())

var _ = bidirekt.Enum(statuses())
