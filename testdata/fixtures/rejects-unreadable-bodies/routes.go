package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Owner struct {
	Name string `json:"name"`
}

var _ = bidirekt.Provides("GET", "/name",
	bidirekt.Response(200, ""),
)

var _ = bidirekt.Provides("GET", "/count",
	bidirekt.Response(200, 0),
)

var _ = bidirekt.Provides("GET", "/anything",
	bidirekt.Response(200, any(nil)),
)

var _ = bidirekt.Provides("GET", "/document",
	bidirekt.Response(200, map[string]any{}),
)

var _ = bidirekt.Provides("GET", "/address",
	bidirekt.Response(200, struct {
		City string `json:"city" bidirekt:"field"`
	}{}),
)

var _ = bidirekt.Provides("GET", "/counts",
	bidirekt.Response(200, map[string]int{}),
)

var _ = bidirekt.Provides("GET", "/owner",
	bidirekt.Response(200, Owner{}),
)

var _ = bidirekt.Provides("GET", "/tags",
	bidirekt.Response(200, []string{}),
)
