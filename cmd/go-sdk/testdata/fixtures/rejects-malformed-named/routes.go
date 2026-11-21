package petstore

import (
	"github.com/bidirekt/go-sdk"
)

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(200, []Pet{}, bidirekt.Named("Pets"), bidirekt.Named("AllPets")),
)

var _ = bidirekt.Provides("GET", "/pets/featured",
	bidirekt.Response(200, []Pet{}, bidirekt.Named("")),
)

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(200, Pet{}, bidirekt.Named("PetResponse")),
)

var _ = bidirekt.Provides("DELETE", "/pets/*",
	bidirekt.Response(204, nil, bidirekt.Named("Nothing")),
)

var _ = bidirekt.Provides("GET", "/pets/all",
	bidirekt.Response(200, Pets{}, bidirekt.Named("AllPets")),
)

var _ = bidirekt.Provides("GET", "/pets/by-name",
	bidirekt.Response(200, PetsByName{}, bidirekt.Named("PetIndex")),
)

var _ = bidirekt.Provides("GET", "/pets/adopted",
	bidirekt.Response(200, []Pet{}, bidirekt.Name{}),
)
