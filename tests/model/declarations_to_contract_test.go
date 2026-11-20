package model_test

import (
	"reflect"
	"testing"

	"github.com/bidirekt/go-sdk/internal/model"
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

	if len(errors) != 0 {
		t.Errorf("errors: got %#v, want none", errors)
	}
	equal(t, "consumes", contract.Consumes, []model.Resource{
		{Provider: "payments_api", Method: "GET", Path: "/charges/*", Responses: []model.Response{{Status: 200}, {Status: 404}}, Location: first},
	})
}

func TestWritesTheSameDerivedListOnceAndRejectsDifferentOnes(t *testing.T) {
	pet := model.NamedSchema{Name: "Pet", Node: model.Node{Shape: model.Object{}}, Location: first}

	same, sameErrors := model.DeclarationsToContract(model.Contract{Schemas: []model.NamedSchema{pet, petList(first, false), petList(second, false)}})
	_, differentErrors := model.DeclarationsToContract(model.Contract{Schemas: []model.NamedSchema{pet, petList(first, false), petList(second, true)}})

	if len(sameErrors) != 0 {
		t.Errorf("errors of the same list: got %#v, want none", sameErrors)
	}
	equal(t, "the same list written once", same.Schemas, []model.NamedSchema{pet, petList(first, false)})
	message := `lists named "PetList" differ: give each list a distinct name`
	equal(t, "errors of different lists", differentErrors, []model.Diagnostic{{Message: message, Location: first}, {Message: message, Location: second}})
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

	equal(t, "errors", errors, []model.Diagnostic{{
		Message:  `schema "Person" is not declared: reference a schema by its name in the contract`,
		Location: second,
	}})
}

func equal(t *testing.T, label string, got, want any) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		return
	}
	if gotText, isText := got.(string); isText {
		t.Errorf("%s:\ngot:\n%s\nwant:\n%s", label, gotText, want)
		return
	}
	t.Errorf("%s:\ngot:  %#v\nwant: %#v", label, got, want)
}
