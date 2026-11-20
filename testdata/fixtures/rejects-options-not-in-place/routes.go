package petstore

import (
	"github.com/bidirekt/go-sdk"
)

var notFound = bidirekt.Response(404, nil)

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(200, Pet{}),
	notFound,
)

var options = []bidirekt.ResourceOption{bidirekt.Response(200, Pet{})}

var _ = bidirekt.Provides("GET", "/pets", options...)

func responses() []bidirekt.ResourceOption {
	return []bidirekt.ResourceOption{bidirekt.Response(200, nil)}
}

var _ = bidirekt.Provides("POST", "/pets",
	bidirekt.Request(NewPet{}),
	bidirekt.Request(Pet{}),
	bidirekt.Response(201, Pet{}),
)
