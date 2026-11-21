package petstore

import (
	"example.com/shared"
)

type Pet struct {
	Name  string       `json:"name" bidirekt:"field"`
	Owner shared.Owner `json:"owner" bidirekt:"field"`
}
