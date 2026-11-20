package petstore

import (
	"example.com/shared"
)

const ContractVersion = shared.Version

type Pet struct {
	Name    string `json:"name" bidirekt:"field"`
	Keeper  any    `json:"keeper" bidirekt:"field type=Owner"`
	Friends []any  `json:"friends" bidirekt:"field type=array items=Friend"`
}
