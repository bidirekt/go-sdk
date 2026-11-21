package reader

import (
	"go/token"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bidirekt/go-sdk/cmd/go-sdk/internal/model"
	"golang.org/x/tools/go/packages"
)

const sdkPath = "github.com/bidirekt/go-sdk"

type loadedPackages struct {
	fset      *token.FileSet
	packages  []*packages.Package
	byPath    map[string]*packages.Package
	inProject map[string]bool
	local     map[string]bool
}

var errorPosition = regexp.MustCompile(`^(.+?):(\d+)(?::\d+)?$`)

var compilerLine = regexp.MustCompile(`^(.+?:\d+(?::\d+)?): (.*)$`)

// The first load lists the import graph without parsing it, to find the packages of local modules (the
// main module, the go.work modules and the replaces by a folder) that the project imports: the second
// load reads their source with the project's, and every other package from export data.
func projectToLoadedPackages(project Project) (loadedPackages, []model.Diagnostic) {
	graph, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedModule,
		Dir:  project.Dir,
	}, "./...")
	if err != nil {
		return loadedPackages{}, []model.Diagnostic{{Message: loadErrorToMessage(err)}}
	}
	inProject := map[string]bool{}
	local := map[string]bool{}
	var patterns []string
	for _, pkg := range graph {
		if isRejectedPattern(pkg) {
			var diagnostics []model.Diagnostic
			for _, packageError := range pkg.Errors {
				diagnostics = append(diagnostics, listErrorToDiagnostics(packageError, project.Dir)...)
			}
			return loadedPackages{}, diagnostics
		}
		inProject[pkg.PkgPath] = true
		patterns = append(patterns, pkg.PkgPath)
	}
	packages.Visit(graph, nil, func(pkg *packages.Package) {
		if !isLocalModule(pkg.Module) {
			return
		}
		local[pkg.PkgPath] = true
		if !inProject[pkg.PkgPath] && !slices.Contains(patterns, pkg.PkgPath) {
			patterns = append(patterns, pkg.PkgPath)
		}
	})
	if len(patterns) == 0 {
		return loadedPackages{inProject: inProject, local: local}, nil
	}
	fset := token.NewFileSet()
	loaded, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes |
			packages.NeedTypesInfo | packages.NeedImports | packages.NeedModule,
		Dir:  project.Dir,
		Fset: fset,
	}, patterns...)
	if err != nil {
		return loadedPackages{}, []model.Diagnostic{{Message: loadErrorToMessage(err)}}
	}
	slices.SortFunc(loaded, func(a, b *packages.Package) int { return strings.Compare(a.PkgPath, b.PkgPath) })
	byPath := map[string]*packages.Package{}
	for _, pkg := range loaded {
		byPath[pkg.PkgPath] = pkg
	}
	return loadedPackages{fset: fset, packages: loaded, byPath: byPath, inProject: inProject, local: local}, nil
}

// When go list rejects the pattern itself, go/packages returns it as a package without name nor module; a
// package whose package clause does not parse has no name either, but has its module.
func isRejectedPattern(pkg *packages.Package) bool {
	return pkg.Name == "" && pkg.Module == nil && len(pkg.Errors) > 0
}

func isLocalModule(module *packages.Module) bool {
	return module != nil && (module.Main || (module.Replace != nil && module.Replace.Version == ""))
}

// The errors of the packages the roots compile with, the roots included. As go build does, a package that
// imports a failing package is not compiled, so its errors, derived from that package's, are not reported.
func packagesToGoErrors(roots []*packages.Package, dir string) []model.Diagnostic {
	failing := map[*packages.Package]bool{}
	visited := map[*packages.Package]bool{}
	var diagnostics []model.Diagnostic
	var visit func(pkg *packages.Package)
	visit = func(pkg *packages.Package) {
		if visited[pkg] {
			return
		}
		visited[pkg] = true
		importsFailing := false
		for _, imported := range pkg.Imports {
			visit(imported)
			importsFailing = importsFailing || failing[imported]
		}
		failing[pkg] = importsFailing || len(pkg.Errors) > 0
		if !importsFailing && len(pkg.Errors) > 0 {
			diagnostics = append(diagnostics, packageToGoErrors(pkg, dir)...)
		}
	}
	for _, root := range roots {
		visit(root)
	}
	return diagnostics
}

// As go build does: the load errors first, then the compiler's output, and go/types only when the compiler
// was not run.
func packageToGoErrors(pkg *packages.Package, dir string) []model.Diagnostic {
	var loadErrors []model.Diagnostic
	var buildOutput []model.Diagnostic
	var sourceErrors []model.Diagnostic
	for _, packageError := range pkg.Errors {
		switch {
		case packageError.Kind == packages.ListError && isBuildOutput(packageError):
			buildOutput = append(buildOutput, listErrorToDiagnostics(packageError, dir)...)
		case packageError.Kind == packages.ListError:
			loadErrors = append(loadErrors, listErrorToDiagnostics(packageError, dir)...)
		default:
			sourceErrors = append(sourceErrors, model.Diagnostic{Message: packageError.Msg, Location: positionToLocation(packageError.Pos, dir)})
		}
	}
	if len(loadErrors) > 0 {
		return loadErrors
	}
	if len(buildOutput) > 0 {
		return buildOutput
	}
	return sourceErrors
}

func isBuildOutput(packageError packages.Error) bool {
	return packageError.Pos == "" && strings.HasPrefix(packageError.Msg, "# ")
}

func listErrorToDiagnostics(packageError packages.Error, dir string) []model.Diagnostic {
	if !isBuildOutput(packageError) {
		return []model.Diagnostic{{Message: packageError.Msg, Location: positionToLocation(packageError.Pos, dir)}}
	}
	var diagnostics []model.Diagnostic
	for _, line := range strings.Split(packageError.Msg, "\n")[1:] {
		if strings.HasPrefix(line, "\t") && len(diagnostics) > 0 {
			diagnostics[len(diagnostics)-1].Message += "\n" + line
			continue
		}
		match := compilerLine.FindStringSubmatch(line)
		if match == nil {
			diagnostics = append(diagnostics, model.Diagnostic{Message: line})
			continue
		}
		diagnostics = append(diagnostics, model.Diagnostic{Message: match[2], Location: positionToLocation(match[1], dir)})
	}
	return diagnostics
}

func positionToLocation(position, dir string) model.SourceLocation {
	match := errorPosition.FindStringSubmatch(position)
	if match == nil {
		return model.SourceLocation{}
	}
	file := match[1]
	if !filepath.IsAbs(file) {
		file = filepath.Join(dir, file)
	}
	line, _ := strconv.Atoi(match[2])
	return model.SourceLocation{File: file, Line: line}
}

func loadErrorToMessage(err error) string {
	message := err.Error()
	if _, stderr, found := strings.Cut(message, "stderr: "); found {
		message = stderr
	}
	return strings.TrimSpace(message)
}
