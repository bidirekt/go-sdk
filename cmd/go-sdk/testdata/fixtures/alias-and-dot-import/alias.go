package petstore

import (
	sdk "github.com/bidirekt/go-sdk"
)

var _ = sdk.Provides("GET", "/pets",
	sdk.Response(200, nil),
)
