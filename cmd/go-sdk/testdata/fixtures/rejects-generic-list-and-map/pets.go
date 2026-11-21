package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type List[T any] []T

type Index[T any] map[string]T

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(200, List[Pet]{}),
)

var _ = bidirekt.Provides("GET", "/owners",
	bidirekt.Response(200, List[Owner]{}),
)

var _ = bidirekt.Provides("GET", "/owners/by-name",
	bidirekt.Response(200, Index[Owner]{}),
)
