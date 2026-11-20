package writer_test

import (
	"testing"

	"github.com/bidirekt/go-sdk/internal/model"
	"github.com/bidirekt/go-sdk/internal/writer"
	"github.com/stretchr/testify/assert"
)

var petLocation = model.SourceLocation{File: "pets/pet.go", Line: 1}

// From the second line on, the expected.yaml of the TypeScript SDK's free-text-quoting fixture: every SDK
// writes the same YAML for the same contract.
const freeTextQuotingYaml = "schemas:\n" +
	"  Pet:\n" +
	"    type: object\n" +
	"    properties:\n" +
	"      plainText:\n" +
	"        type: string\n" +
	"        description: plain text\n" +
	"      looksLikeYes:\n" +
	"        type: string\n" +
	"        description: yes\n" +
	"      looksLikeOn:\n" +
	"        type: string\n" +
	"        description: on\n" +
	"      looksLikeTrue:\n" +
	"        type: string\n" +
	"        description: \"true\"\n" +
	"      looksLikeNull:\n" +
	"        type: string\n" +
	"        description: \"null\"\n" +
	"      \"~\":\n" +
	"        type: string\n" +
	"      looksLikeInteger:\n" +
	"        type: string\n" +
	"        description: \"123\"\n" +
	"      looksLikeHexadecimal:\n" +
	"        type: string\n" +
	"        description: \"0x1F\"\n" +
	"      2001-12-14:\n" +
	"        type: string\n" +
	"      \"a: b\":\n" +
	"        type: string\n" +
	"      a:b:\n" +
	"        type: string\n" +
	"      \"a #b\":\n" +
	"        type: string\n" +
	"      a#b:\n" +
	"        type: string\n" +
	"      \"#x\":\n" +
	"        type: string\n" +
	"      \"*x\":\n" +
	"        type: string\n" +
	"      \"&x\":\n" +
	"        type: string\n" +
	"      \"!x\":\n" +
	"        type: string\n" +
	"      \"|x\":\n" +
	"        type: string\n" +
	"      \">x\":\n" +
	"        type: string\n" +
	"      \"'x\":\n" +
	"        type: string\n" +
	"      '\"x':\n" +
	"        type: string\n" +
	"      \"%x\":\n" +
	"        type: string\n" +
	"      \"@x\":\n" +
	"        type: string\n" +
	"      \"`x\":\n" +
	"        type: string\n" +
	"      \"- x\":\n" +
	"        type: string\n" +
	"      \"? x\":\n" +
	"        type: string\n" +
	"      \",x\":\n" +
	"        type: string\n" +
	"      -x:\n" +
	"        type: string\n" +
	"      ?x:\n" +
	"        type: string\n" +
	"      :x:\n" +
	"        type: string\n" +
	"      \" x\":\n" +
	"        type: string\n" +
	"      \"x \":\n" +
	"        type: string\n" +
	"      it's:\n" +
	"        type: string\n" +
	"      say \"hi\":\n" +
	"        type: string\n" +
	"      a\\b:\n" +
	"        type: string\n"

func TestQuotesDescriptionsAndPropertyNamesLikeTheTypeScriptSDK(t *testing.T) {
	properties := []model.Property{
		{Name: "plainText", Node: model.Node{Shape: model.Primitive{Type: "string"}, Description: "plain text"}, Location: petLocation},
		{Name: "looksLikeYes", Node: model.Node{Shape: model.Primitive{Type: "string"}, Description: "yes"}, Location: petLocation},
		{Name: "looksLikeOn", Node: model.Node{Shape: model.Primitive{Type: "string"}, Description: "on"}, Location: petLocation},
		{Name: "looksLikeTrue", Node: model.Node{Shape: model.Primitive{Type: "string"}, Description: "true"}, Location: petLocation},
		{Name: "looksLikeNull", Node: model.Node{Shape: model.Primitive{Type: "string"}, Description: "null"}, Location: petLocation},
		{Name: "~", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "looksLikeInteger", Node: model.Node{Shape: model.Primitive{Type: "string"}, Description: "123"}, Location: petLocation},
		{Name: "looksLikeHexadecimal", Node: model.Node{Shape: model.Primitive{Type: "string"}, Description: "0x1F"}, Location: petLocation},
		{Name: "2001-12-14", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "a: b", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "a:b", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "a #b", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "a#b", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "#x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "*x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "&x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "!x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "|x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: ">x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "'x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "\"x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "%x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "@x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "`x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "- x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "? x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: ",x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "-x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "?x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: ":x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: " x", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "x ", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "it's", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "say \"hi\"", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
		{Name: "a\\b", Node: model.Node{Shape: model.Primitive{Type: "string"}}, Location: petLocation},
	}
	contract := model.Contract{Schemas: []model.NamedSchema{
		{Name: "Pet", Node: model.Node{Shape: model.Object{Properties: properties}}, Location: petLocation},
	}}

	assert.Equal(t, "# Generated by github.com/bidirekt/go-sdk. Do not edit.\n"+freeTextQuotingYaml, writer.ContractToYaml(contract))
}
