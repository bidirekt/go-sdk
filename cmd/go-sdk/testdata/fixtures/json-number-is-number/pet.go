package petstore

import (
	"encoding/json"
)

type Pet struct {
	Weight json.Number `json:"weight" bidirekt:"field"`
}
