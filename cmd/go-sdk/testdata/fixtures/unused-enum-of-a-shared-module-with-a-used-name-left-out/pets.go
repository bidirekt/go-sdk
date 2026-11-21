package petstore

import (
	"example.com/shared"
	"github.com/bidirekt/go-sdk"
)

type Status string

const (
	StatusAvailable Status = "available"
	StatusAdopted   Status = "adopted"
)

var _ = bidirekt.Enum(StatusAvailable, StatusAdopted)

type Pet struct {
	Status Status       `json:"status" bidirekt:"field"`
	Owner  shared.Owner `json:"owner" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(200, Pet{}),
)
