package pets

import (
	"time"

	"github.com/bidirekt/go-sdk"
)

type Status string

const (
	StatusAvailable Status = "available"
	StatusPending   Status = "pending"
	StatusAdopted   Status = "adopted"
)

var _ = bidirekt.Enum(StatusAvailable, StatusPending, StatusAdopted)

type Pet struct {
	ID        int64     `json:"id" bidirekt:"field"`
	Name      string    `json:"name" bidirekt:"field"`
	Status    Status    `json:"status" bidirekt:"field"`
	Tags      []string  `json:"tags" bidirekt:"field"`
	Weight    float64   `json:"weight" bidirekt:"field description='In kilograms'"`
	BirthDate *string   `json:"birthDate" bidirekt:"field"`
	CreatedAt time.Time `json:"createdAt" bidirekt:"field type=string"`
}

type NewPet struct {
	Name      string   `json:"name" bidirekt:"field"`
	Tags      []string `json:"tags,omitempty" bidirekt:"field"`
	Weight    float64  `json:"weight" bidirekt:"field"`
	BirthDate *string  `json:"birthDate,omitempty" bidirekt:"field"`
}

type StatusChange struct {
	Status Status `json:"status" bidirekt:"field"`
}
