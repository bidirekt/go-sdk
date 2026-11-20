package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(200, []Pet{}),
)

var _ = bidirekt.Provides("GET", "/pets/featured",
	bidirekt.Response(200, []*Pet{}),
)

var _ = bidirekt.Provides("GET", "/residents",
	bidirekt.Response(200, []Pet{}, bidirekt.Named("Residents")),
)

var _ = bidirekt.Provides("GET", "/visitors",
	bidirekt.Response(200, []string{}, bidirekt.Named("Residents")),
)
