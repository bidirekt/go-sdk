package main

import (
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Consumes("petstore_api", "GET", "/pets",
	bidirekt.Response(200, nil),
)

func main() {}
