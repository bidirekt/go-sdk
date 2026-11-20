package petstore

import (
	"time"

	"github.com/bidirekt/go-sdk"
)

type Owner struct {
	Name string `json:"name"`
}

type Pet struct {
	Name string `json:"name" bidirekt:"field"`
}

type Price struct {
	Cents int
}

func (Price) MarshalJSON() ([]byte, error) {
	return []byte(`"0.00"`), nil
}

type Label struct {
	Text string
}

func (*Label) MarshalText() ([]byte, error) {
	return []byte("label"), nil
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

var _ = bidirekt.Provides("GET", "/price",
	bidirekt.Response(200, Price{}),
)

var _ = bidirekt.Provides("GET", "/label",
	bidirekt.Response(200, &Label{}),
)

var _ = bidirekt.Provides("GET", "/dates",
	bidirekt.Response(200, []time.Time{}, bidirekt.Named("Dates")),
)

var _ = bidirekt.Provides("GET", "/prices",
	bidirekt.Response(200, map[string]Price{}, bidirekt.Named("Prices")),
)

var _ = bidirekt.Provides("GET", "/pets-by-id",
	bidirekt.Response(200, map[int]Pet{}, bidirekt.Named("PetsByID")),
)
