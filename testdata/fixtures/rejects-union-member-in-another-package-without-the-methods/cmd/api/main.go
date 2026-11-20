package main

import (
	"example.com/petstore/pets"
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(200, pets.Pet{}),
)

func main() {}
