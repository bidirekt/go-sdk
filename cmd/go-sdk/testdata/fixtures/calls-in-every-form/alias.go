package petstore

import (
	sdk "github.com/bidirekt/go-sdk"
)

type PetSize string

var _ = (sdk.Enum[PetSize])("small", "large")

var _ = (sdk.Provides)("GET", "/pets",
	(sdk.Response)(200, nil),
)
