package shared

import (
	"github.com/bidirekt/go-sdk"
)

type LineItem struct {
	Description string `json:"description" bidirekt:"field"`
	Cents       int    `json:"cents" bidirekt:"field"`
}

type Invoice struct {
	Number string     `json:"number" bidirekt:"field"`
	Items  []LineItem `json:"items" bidirekt:"field"`
}

var _ = bidirekt.Schema(Invoice{}, bidirekt.Named("Bill"))

type Customer struct {
	Name string `json:"name" bidirekt:"field"`
}
