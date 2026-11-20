package petstore

import (
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("get", "/pets",
	bidirekt.Response(200, nil),
)

var _ = bidirekt.Consumes("petstore_api", "Post", "/pets",
	bidirekt.Response(201, nil),
)
