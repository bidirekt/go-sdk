package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Status string

const (
	StatusAvailable Status = "available"
	StatusSold      Status = "sold"
)

type Size string

const (
	SizeSmall Size = "small"
	SizeLarge Size = "large"
)

var _ = bidirekt.Enum(StatusAvailable, StatusSold)

var _ = bidirekt.Schema(Status(""), bidirekt.Named("PetStatus"))

var _ = bidirekt.Enum(SizeSmall, SizeLarge)
