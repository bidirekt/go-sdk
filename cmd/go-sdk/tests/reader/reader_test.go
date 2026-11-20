package reader_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bidirekt/go-sdk/cmd/go-sdk/internal/model"
	"github.com/bidirekt/go-sdk/cmd/go-sdk/internal/reader"
)

// The cases are read in place, so the caller's go.work or -mod=mod must not reach them.
func TestMain(m *testing.M) {
	for key, value := range map[string]string{"GOWORK": "off", "GOFLAGS": "-mod=readonly"} {
		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}
	}
	os.Exit(m.Run())
}

func caseDir(t *testing.T, name string) string {
	t.Helper()
	dir, err := filepath.Abs(filepath.Join("..", "..", "testdata", "reader", name))
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}
	return dir
}

func at(dir, file string, line int) model.SourceLocation {
	return model.SourceLocation{File: filepath.Join(dir, file), Line: line}
}

func TestRecognizesCallsByTheirPackageOfOrigin(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "calls")

	result := reader.ProjectToResult(reader.Project{Dir: dir})

	if len(result.Errors) != 0 {
		t.Errorf("errors: got %#v, want none", result.Errors)
	}
	equal(t, "an import with an alias, a call inside a function, http.MethodGet and http.StatusOK count; a Provides of another package, a _test.go file and a file outside the build tags do not", result.Contract.Provides, []model.Resource{
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
	})
	equal(t, "a dot import counts", result.Contract.Consumes, []model.Resource{
		{
			Provider:  "petstore_reviews",
			Method:    "POST",
			Path:      "/dot",
			Request:   &model.Body{Location: at(dir, "dot.go", 8)},
			Responses: []model.Response{{Status: 201, Body: model.Body{Location: at(dir, "dot.go", 9)}}},
			Location:  at(dir, "dot.go", 7),
		},
	})
}

func TestReadsTheUsedEnumsButNotResourcesOfAnotherWorkspaceModule(t *testing.T) {
	dir := caseDir(t, "workspace")
	t.Setenv("GOWORK", filepath.Join(dir, "go.work"))

	result := reader.ProjectToResult(reader.Project{Dir: filepath.Join(dir, "api")})

	if len(result.Errors) != 0 {
		t.Errorf("errors: got %#v, want none", result.Errors)
	}
	equal(t, "resources of the project only", result.Contract.Provides, []model.Resource{
		{
			Method:    "GET",
			Path:      "/pets",
			Responses: []model.Response{{Status: 200, Body: model.Body{Schema: "Pet", Location: at(dir, "api/routes.go", 13)}}},
			Location:  at(dir, "api/routes.go", 12),
		},
	})
	equal(t, "the enum of the shared module the contract uses, and not the one it does not", result.Contract.Schemas, []model.NamedSchema{
		{
			Name: "Pet",
			Node: model.Node{Shape: model.Object{Properties: []model.Property{
				{Name: "status", Node: model.Node{Shape: model.Ref{Schema: "PetStatus"}}, Location: at(dir, "api/routes.go", 9)},
			}}},
			Location: at(dir, "api/routes.go", 8),
		},
		{
			Name:     "PetStatus",
			Node:     model.Node{Shape: model.Enum{Members: []string{"available", "adopted"}}},
			Location: at(dir, "shared/status.go", 14),
		},
	})
}

func TestReportsGoMessageWhenGoRejectsTheProjectFolder(t *testing.T) {
	dir := caseDir(t, "outside-workspace")
	t.Setenv("GOWORK", filepath.Join(dir, "go.work"))

	result := reader.ProjectToResult(reader.Project{Dir: filepath.Join(dir, "worker")})

	equal(t, "a folder outside the modules of the go.work", result.Errors, []model.Diagnostic{
		{Message: "pattern ./...: directory prefix . does not contain modules listed in go.work or their selected dependencies"},
	})
}

func TestReadsTheStructsATagNamesFromADependency(t *testing.T) {
	dir := caseDir(t, "tag-names-structs-of-a-dependency")
	t.Setenv("GOFLAGS", "-mod=vendor")

	result := reader.ProjectToResult(reader.Project{Dir: dir})

	if len(result.Errors) != 0 {
		t.Errorf("errors: got %#v, want none", result.Errors)
	}
	properties := map[string][]string{}
	for _, schema := range result.Contract.Schemas {
		for _, property := range schema.Node.Shape.(model.Object).Properties {
			properties[schema.Name] = append(properties[schema.Name], property.Name)
		}
	}
	equal(t, "the structs named by type= and items=, read from the export data of a vendored module", properties, map[string][]string{
		"Friend": {"name"},
		"Owner":  {"name"},
		"Pet":    {"name", "keeper", "friends"},
	})
}

func TestRejectsANonConstantArgumentAtItsLine(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "non-constant")

	result := reader.ProjectToResult(reader.Project{Dir: dir})

	equal(t, "errors", result.Errors, []model.Diagnostic{
		{
			Message:  `"petsPath" cannot be read without running the code: write a literal or a constant`,
			Location: at(dir, "routes.go", 9),
		},
	})
}

func TestReportsACompileErrorInAPackageWithADeclarationWithGoMessage(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "compile-error")

	result := reader.ProjectToResult(reader.Project{Dir: dir})

	equal(t, "errors", result.Errors, []model.Diagnostic{
		{
			Message:  `cannot use "many" (untyped string constant) as int value in variable declaration`,
			Location: at(dir, "routes.go", 7),
		},
	})
}

func TestReportsAMissingGoSumEntryWithGoMessage(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "missing-go-sum")

	result := reader.ProjectToResult(reader.Project{Dir: dir})

	if len(result.Errors) != 1 {
		t.Fatalf("errors: got %d, want 1: %#v", len(result.Errors), result.Errors)
	}
	if prefix := "missing go.sum entry for module providing package golang.org/x/sync/errgroup"; !strings.HasPrefix(result.Errors[0].Message, prefix) {
		t.Errorf("message: got %q, want it to start with %q", result.Errors[0].Message, prefix)
	}
	equal(t, "location", result.Errors[0].Location, at(dir, "routes.go", 5))
}

func TestIgnoresAnErrorInAPackageWithoutDeclarations(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "error-outside-declarations")

	result := reader.ProjectToResult(reader.Project{Dir: dir})

	if len(result.Errors) != 0 {
		t.Errorf("errors: got %#v, want none", result.Errors)
	}
	equal(t, "provides", result.Contract.Provides, []model.Resource{
		{
			Method:    "GET",
			Path:      "/pets",
			Responses: []model.Response{{Status: 200, Body: model.Body{Location: at(dir, "routes.go", 8)}}},
			Location:  at(dir, "routes.go", 7),
		},
	})
}

func TestFindsTheModuleOfTheNearestGoModAtOrAboveTheFolder(t *testing.T) {
	t.Parallel()
	dir := caseDir(t, "calls")

	project, found := reader.DirToProject(filepath.Join(dir, "fake"))
	_, foundAboveRoot := reader.DirToProject(t.TempDir())

	if !found {
		t.Errorf("a folder below go.mod: got false, want true")
	}
	equal(t, "the project folder stays the one given", project.Dir, filepath.Join(dir, "fake"))
	if foundAboveRoot {
		t.Errorf("a folder without go.mod at or above it: got true, want false")
	}
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
