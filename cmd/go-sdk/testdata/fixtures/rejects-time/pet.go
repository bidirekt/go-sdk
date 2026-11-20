package petstore

import (
	"time"
)

type Pet struct {
	Name   string    `json:"name" bidirekt:"field"`
	BornAt time.Time `json:"bornAt" bidirekt:"field"`
}
