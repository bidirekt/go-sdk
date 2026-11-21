package pets

import "github.com/bidirekt/go-sdk"

type Status string

const (
	StatusAvailable Status = "available"
	StatusPending   Status = "pending"
	StatusSold      Status = "sold"
)

var _ = bidirekt.Enum(StatusAvailable, StatusPending, StatusSold)

type Pet struct {
	ID        int64   `json:"id" bidirekt:"field"`
	Name      string  `json:"name" bidirekt:"field"`
	Status    Status  `json:"status" bidirekt:"field"`
	BirthDate *string `json:"birthDate" bidirekt:"field description='ISO 8601 date'"`
}

type NewPet struct {
	Name      string  `json:"name" bidirekt:"field"`
	BirthDate *string `json:"birthDate,omitempty" bidirekt:"field"`
}
