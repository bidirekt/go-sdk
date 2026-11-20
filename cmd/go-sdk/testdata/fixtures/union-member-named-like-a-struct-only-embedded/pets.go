package app

import (
	"example.com/app/users"
	"github.com/bidirekt/go-sdk"
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
	AdoptedFrom Adopter `json:"adoptedFrom" bidirekt:"field union={person: Owner, shelter: Shelter}"`
}

var _ = bidirekt.Provides("GET", "/pets", bidirekt.Response(200, Pet{}))
var _ = bidirekt.Provides("GET", "/accounts", bidirekt.Response(200, users.Account{}))
