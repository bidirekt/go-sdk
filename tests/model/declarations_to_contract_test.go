package model_test

import (
	"testing"

	"github.com/bidirekt/go-sdk/internal/model"
	"github.com/stretchr/testify/assert"
)

var (
	first  = model.SourceLocation{File: "/petstore/clients.go", Line: 10}
	second = model.SourceLocation{File: "/petstore/refunds.go", Line: 20}
)

func petList(location model.SourceLocation, nullable bool) model.NamedSchema {
	items := model.Node{Shape: model.Ref{Schema: "Pet"}, Nullable: nullable}
	return model.NamedSchema{Name: "PetList", Node: model.Node{Shape: model.Array{Items: items}}, Location: location, Derived: true}
}

func TestJoinsTheStatusesOfAResourceConsumedInTwoDeclarations(t *testing.T) {
	declarations := model.Contract{Consumes: []model.Resource{
		{Provider: "payments_api", Method: "GET", Path: "/charges/*", Responses: []model.Response{{Status: 200}}, Location: first},
		{Provider: "payments_api", Method: "GET", Path: "/charges/*", Responses: []model.Response{{Status: 200}, {Status: 404}}, Location: second},
	}}

	contract, errors := model.DeclarationsToContract(declarations)

	assert.Empty(t, errors, "errors")
	assert.Equal(t, []model.Resource{
		{Provider: "payments_api", Method: "GET", Path: "/charges/*", Responses: []model.Response{{Status: 200}, {Status: 404}}, Location: first},
	}, contract.Consumes)
}

func TestWritesTheSameDerivedListOnceAndRejectsDifferentOnes(t *testing.T) {
	pet := model.NamedSchema{Name: "Pet", Node: model.Node{Shape: model.Object{}}, Location: first}

	same, sameErrors := model.DeclarationsToContract(model.Contract{Schemas: []model.NamedSchema{pet, petList(first, false), petList(second, false)}})
	_, differentErrors := model.DeclarationsToContract(model.Contract{Schemas: []model.NamedSchema{pet, petList(first, false), petList(second, true)}})

	assert.Empty(t, sameErrors, "errors of the same list")
	assert.Equal(t, []model.NamedSchema{pet, petList(first, false)}, same.Schemas, "the same list written once")
	message := `lists named "PetList" differ: give each list a distinct name`
	assert.Equal(t, []model.Diagnostic{{Message: message, Location: first}, {Message: message, Location: second}}, differentErrors, "errors of different lists")
}

func TestRejectsAReferenceToASchemaNotDeclared(t *testing.T) {
	declarations := model.Contract{Schemas: []model.NamedSchema{{
		Name: "Pet",
		Node: model.Node{Shape: model.Object{Properties: []model.Property{
			{Name: "owner", Node: model.Node{Shape: model.Ref{Schema: "Person"}}, Location: second},
		}}},
		Location: first,
	}}}

	_, errors := model.DeclarationsToContract(declarations)

	assert.Equal(t, []model.Diagnostic{{
		Message:  `schema "Person" is not declared: reference a schema by its name in the contract`,
		Location: second,
	}}, errors)
}
