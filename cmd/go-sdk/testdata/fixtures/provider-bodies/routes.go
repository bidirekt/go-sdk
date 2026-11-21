package petstore

import (
	"net/http"

	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides(http.MethodGet, "/pets/*",
	bidirekt.Response(http.StatusOK, Pet{}),
	bidirekt.Response(http.StatusNotFound, nil),
)

var _ = bidirekt.Provides(http.MethodPut, "/pets/*",
	bidirekt.Request(&NewPet{}),
	bidirekt.Response(http.StatusOK, &Pet{}),
)

var _ = bidirekt.Provides(http.MethodGet, "/pets",
	bidirekt.Response(http.StatusOK, []Pet{}),
)

var _ = bidirekt.Provides(http.MethodGet, "/pets/featured",
	bidirekt.Response(http.StatusOK, []*Pet{}, bidirekt.Named("FeaturedPets")),
)

var _ = bidirekt.Provides(http.MethodGet, "/pets/adopted",
	bidirekt.Response(http.StatusOK, Pets{}),
)

var _ = bidirekt.Provides(http.MethodGet, "/pets/by-name",
	bidirekt.Response(http.StatusOK, PetsByName{}),
)

var _ = bidirekt.Provides(http.MethodGet, "/pets/counts",
	bidirekt.Response(http.StatusOK, map[string]int{}, bidirekt.Named("PetCounts")),
)

var _ = bidirekt.Provides(http.MethodPost, "/pets/*/adoptions",
	bidirekt.Request(struct {
		Owner string `json:"owner" bidirekt:"field"`
	}{}, bidirekt.Named("Adoption")),
	bidirekt.Response(http.StatusCreated, nil),
)

var _ = bidirekt.Provides(http.MethodDelete, "/pets/*",
	bidirekt.Request(nil),
	bidirekt.Response(http.StatusNoContent, nil),
)
