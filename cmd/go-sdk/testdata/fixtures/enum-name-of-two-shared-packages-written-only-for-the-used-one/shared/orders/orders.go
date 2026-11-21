package orders

import "github.com/bidirekt/go-sdk"

type Status string

const (
	StatusPlaced  Status = "placed"
	StatusShipped Status = "shipped"
)

var _ = bidirekt.Enum(StatusPlaced, StatusShipped)
