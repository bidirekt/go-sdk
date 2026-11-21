package billing

import "github.com/bidirekt/go-sdk"

type Status string

const (
	StatusDue  Status = "due"
	StatusPaid Status = "paid"
)

var _ = bidirekt.Enum(StatusDue, StatusPaid)

type Money struct {
	Cents int64 `json:"cents" bidirekt:"field"`
}
