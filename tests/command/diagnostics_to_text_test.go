package command_test

import (
	"testing"

	"github.com/bidirekt/go-sdk/internal/command"
	"github.com/bidirekt/go-sdk/internal/model"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, `api/routes.go:30 error in another file
pets.go:3 warning: a near miss
pets.go:9 second error
`, text)
}

func TestShowsADiagnosticWithoutLineOrPlace(t *testing.T) {
	errors := []model.Diagnostic{
		{Message: "go: errors parsing go.mod"},
		{Message: "a whole file", Location: model.SourceLocation{File: "/home/dev/shared/status.go"}},
	}

	text := command.DiagnosticsToText(errors, nil, cwd)

	assert.Equal(t, "go: errors parsing go.mod\n../shared/status.go a whole file\n", text)
}
