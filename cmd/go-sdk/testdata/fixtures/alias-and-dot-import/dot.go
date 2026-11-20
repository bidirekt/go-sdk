package petstore

import (
	. "github.com/bidirekt/go-sdk"
)

var _ = Provides("POST", "/pets",
	Request(nil),
	Response(201, nil),
)
