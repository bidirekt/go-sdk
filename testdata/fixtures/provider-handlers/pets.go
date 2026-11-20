package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

type NewPet struct {
	Name string `json:"name" bidirekt:"field"`
}

type PetCard struct {
	Name string `json:"name" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(200, Pet{}),
	bidirekt.Response(404, nil),
)

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(200, []Pet{}),
)

var _ = bidirekt.Provides("GET", "/pets/featured",
	bidirekt.Response(200, []Pet{}, bidirekt.Named("PetCollection")),
)

var _ = bidirekt.Provides("GET", "/pets/*/card",
	bidirekt.Response(200, PetCard{}),
)

var _ = bidirekt.Provides("POST", "/pets",
	bidirekt.Request(NewPet{}),
	bidirekt.Response(201, Pet{}),
	bidirekt.Response(400, nil),
	bidirekt.Response(409, nil),
)

var _ = bidirekt.Provides("PUT", "/pets/*",
	bidirekt.Request(nil),
	bidirekt.Response(204, nil),
)
