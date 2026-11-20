package petstore

import (
	"github.com/bidirekt/go-sdk"
)

type Tree []Tree

type Categories map[string]Categories

type Lineage []*Lineage

type Pet struct {
	Name       string     `json:"name" bidirekt:"field"`
	Family     Tree       `json:"family" bidirekt:"field"`
	Categories Categories `json:"categories" bidirekt:"field"`
	Lineage    Lineage    `json:"lineage" bidirekt:"field"`
	Breeders   Tree       `json:"breeders" bidirekt:"field union={kennel: Pet}"`
}

var _ = bidirekt.Provides("GET", "/trees",
	bidirekt.Response(200, Tree{}),
)
