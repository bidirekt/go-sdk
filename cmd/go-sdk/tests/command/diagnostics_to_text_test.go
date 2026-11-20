package command_test

import (
	"testing"

	"github.com/bidirekt/go-sdk/cmd/go-sdk/internal/command"
	"github.com/bidirekt/go-sdk/cmd/go-sdk/internal/model"
)

const cwd = "/home/dev/petstore"

func TestSortsErrorsAndWarningsTogetherByFileAndLineWithoutRepeats(t *testing.T) {
	errors := []model.Diagnostic{
		{Message: "second error", Location: model.SourceLocation{File: "/home/dev/petstore/pets.go", Line: 9}},
		{Message: "error in another file", Location: model.SourceLocation{File: "/home/dev/petstore/api/routes.go", Line: 30}},
		{Message: "second error", Location: model.SourceLocation{File: "/home/dev/petstore/pets.go", Line: 9}},
	}
	warnings := []model.Diagnostic{
		{Message: "a near miss", Location: model.SourceLocation{File: "/home/dev/petstore/pets.go", Line: 3}},
	}

	text := command.DiagnosticsToText(errors, warnings, cwd)

	want := `api/routes.go:30 error in another file
pets.go:3 warning: a near miss
pets.go:9 second error
`
	if text != want {
		t.Errorf("text:\ngot:\n%s\nwant:\n%s", text, want)
	}
}

func TestShowsADiagnosticWithoutPlace(t *testing.T) {
	errors := []model.Diagnostic{
		{Message: "go: errors parsing go.mod"},
	}

	text := command.DiagnosticsToText(errors, nil, cwd)

	if want := "go: errors parsing go.mod\n"; text != want {
		t.Errorf("text:\ngot:\n%s\nwant:\n%s", text, want)
	}
}
