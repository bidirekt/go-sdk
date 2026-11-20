package bidirekt_test

import (
	"net/http"
	"os/exec"
	"strings"
	"testing"

	"github.com/bidirekt/go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusAdopted   PetStatus = "adopted"
)

type Adopter interface{ adopter() }

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Shelter struct {
	City string `json:"city" bidirekt:"field"`
}

func (Owner) adopter()   {}
func (Shelter) adopter() {}

type Pet struct {
	ID          int64     `json:"id" bidirekt:"field"`
	Status      PetStatus `json:"status" bidirekt:"field"`
	AdoptedFrom Adopter   `json:"adoptedFrom" bidirekt:"field union={person: Owner, shelter: Shelter}"`
}

type NewPet struct {
	Name string `json:"name" bidirekt:"field"`
}

type Rating struct {
	Stars int `json:"stars" bidirekt:"field"`
}

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(http.StatusOK, Pet{}),
	bidirekt.Response(http.StatusNotFound, nil),
)

var _ = bidirekt.Provides(http.MethodPost, "/pets",
	bidirekt.Request(NewPet{}),
	bidirekt.Response(http.StatusCreated, Pet{}),
	bidirekt.Response(http.StatusBadRequest, nil),
)

var _ = bidirekt.Provides("GET", "/pets",
	bidirekt.Response(http.StatusOK, []Pet{}),
)

var _ = bidirekt.Provides("GET", "/pets/adopted",
	bidirekt.Response(http.StatusOK, []Pet{}, bidirekt.Named("AdoptedPets")),
)

var _ = bidirekt.Consumes("petstore_reviews", "GET", "/ratings/*",
	bidirekt.Response(http.StatusOK, Rating{}),
	bidirekt.Response(http.StatusNotFound, nil),
)

var _ = bidirekt.Enum(PetStatusAvailable, PetStatusAdopted)

var _ = bidirekt.Enum[PetStatus]("available", "adopted")

var _ = bidirekt.Schema(Pet{}, bidirekt.Named("PetResponse"))

var _ = bidirekt.Schema(PetStatus(""), bidirekt.Named("Status"))

func TestDeclarationsCompileAsPackageVariables(t *testing.T) {
	resource := bidirekt.Provides("DELETE", "/pets/*", bidirekt.Response(http.StatusNoContent, nil))
	schema := bidirekt.Schema(Owner{}, bidirekt.Named("PetOwner"))
	enum := bidirekt.Enum(PetStatusAvailable)

	assert.Equal(t, bidirekt.Resource{}, resource, "Provides does nothing at runtime")
	assert.Equal(t, bidirekt.SchemaDeclaration{}, schema, "Schema does nothing at runtime")
	assert.Equal(t, bidirekt.EnumDeclaration{}, enum, "Enum does nothing at runtime")
}

func TestPackageImportsNothing(t *testing.T) {
	output, err := exec.Command("go", "list", "-deps", "github.com/bidirekt/go-sdk").Output()
	require.NoError(t, err, "go list")

	assert.Equal(t, []string{"github.com/bidirekt/go-sdk"}, strings.Fields(string(output)), "dependencies of the bidirekt package")
}
