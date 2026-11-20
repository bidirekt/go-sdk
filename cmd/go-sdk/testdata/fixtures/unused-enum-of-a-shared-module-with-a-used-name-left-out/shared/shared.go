package shared

import "github.com/bidirekt/go-sdk"

type Status string

const (
	StatusOpen   Status = "open"
	StatusClosed Status = "closed"
)

var _ = bidirekt.Enum(StatusOpen, StatusClosed)

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}
