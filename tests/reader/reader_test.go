package reader_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/bidirekt/go-sdk/internal/model"
	"github.com/bidirekt/go-sdk/internal/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func caseDir(t *testing.T, name string) string {
	t.Helper()
	dir, err := filepath.Abs(filepath.Join("..", "..", "testdata", "reader", name))
	require.NoError(t, err)
	return dir
}

func at(dir, file string, line int) model.SourceLocation {
	return model.SourceLocation{File: filepath.Join(dir, file), Line: line}
}

func TestRecognizesCallsByTheirPackageOfOrigin(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "calls")

	result := reader.ProjectToResult(reader.Project{Dir: dir})

	assert.Empty(t, result.Errors, "errors")
	assert.Equal(t, []model.Resource{
		{
			Method:    "GET",
			Path:      "/alias",
			Responses: []model.Response{{Status: 200, Body: model.Body{Location: at(dir, "alias.go", 10)}}},
			Location:  at(dir, "alias.go", 9),
		},
		{
			Method:    "PUT",
			Path:      "/inside",
			Responses: []model.Response{{Status: 204, Body: model.Body{Location: at(dir, "function.go", 9)}}},
			Location:  at(dir, "function.go", 8),
		},
	}, result.Contract.Provides, "an import with an alias, a call inside a function, http.MethodGet and http.StatusOK count; a Provides of another package, a _test.go file and a file outside the build tags do not")
	assert.Equal(t, []model.Resource{
		{
			Provider:  "petstore_reviews",
			Method:    "POST",
			Path:      "/dot",
			Request:   &model.Body{Location: at(dir, "dot.go", 8)},
			Responses: []model.Response{{Status: 201, Body: model.Body{Location: at(dir, "dot.go", 9)}}},
			Location:  at(dir, "dot.go", 7),
		},
	}, result.Contract.Consumes, "a dot import counts")
}

func TestReadsEnumsButNotResourcesOfAnotherWorkspaceModule(t *testing.T) {
	dir := caseDir(t, "workspace")
	t.Setenv("GOWORK", filepath.Join(dir, "go.work"))

	result := reader.ProjectToResult(reader.Project{Dir: filepath.Join(dir, "api")})

	assert.Empty(t, result.Errors, "errors")
	assert.Equal(t, []model.Resource{
		{
			Method:    "GET",
			Path:      "/pets",
			Responses: []model.Response{{Status: 200, Body: model.Body{Location: at(dir, "api/routes.go", 11)}}},
			Location:  at(dir, "api/routes.go", 10),
		},
	}, result.Contract.Provides, "resources of the project only")
	assert.Equal(t, []model.NamedSchema{
		{
			Name:     "PetStatus",
			Node:     model.Node{Shape: model.Enum{Members: []string{"available", "adopted"}}},
			Location: at(dir, "shared/status.go", 14),
		},
	}, result.Contract.Schemas, "the enum of the shared module")
}

func TestRejectsANonConstantArgumentAtItsLine(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "non-constant")

	result := reader.ProjectToResult(reader.Project{Dir: dir})

	assert.Equal(t, []model.Diagnostic{
		{
			Message:  `"petsPath" cannot be read without running the code: write a literal or a constant`,
			Location: at(dir, "routes.go", 9),
		},
	}, result.Errors)
}

func TestReportsACompileErrorInAPackageWithADeclarationWithGoMessage(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "compile-error")

	result := reader.ProjectToResult(reader.Project{Dir: dir})

	assert.Equal(t, []model.Diagnostic{
		{
			Message:  `cannot use "many" (untyped string constant) as int value in variable declaration`,
			Location: at(dir, "routes.go", 7),
		},
	}, result.Errors)
}

func TestReportsAMissingGoSumEntryWithGoMessage(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "missing-go-sum")

	result := reader.ProjectToResult(reader.Project{Dir: dir})

	require.Len(t, result.Errors, 1, "errors")
	assert.True(t, strings.HasPrefix(result.Errors[0].Message, "missing go.sum entry for module providing package golang.org/x/sync/errgroup"), "message %q", result.Errors[0].Message)
	assert.Equal(t, at(dir, "routes.go", 5), result.Errors[0].Location, "location")
}

func TestIgnoresAnErrorInAPackageWithoutDeclarations(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "error-outside-declarations")

	result := reader.ProjectToResult(reader.Project{Dir: dir})

	assert.Empty(t, result.Errors, "errors")
	assert.Equal(t, []model.Resource{
		{
			Method:    "GET",
			Path:      "/pets",
			Responses: []model.Response{{Status: 200, Body: model.Body{Location: at(dir, "routes.go", 8)}}},
			Location:  at(dir, "routes.go", 7),
		},
	}, result.Contract.Provides)
}

func TestFindsTheModuleOfTheNearestGoModAtOrAboveTheFolder(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "calls")

	project, found := reader.DirToProject(filepath.Join(dir, "fake"))
	_, foundAboveRoot := reader.DirToProject(t.TempDir())

	assert.True(t, found, "a folder below go.mod")
	assert.Equal(t, filepath.Join(dir, "fake"), project.Dir, "the project folder stays the one given")
	assert.False(t, foundAboveRoot, "a folder without go.mod at or above it")
}
