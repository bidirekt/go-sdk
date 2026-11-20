package reader

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"github.com/bidirekt/go-sdk/internal/model"
	"golang.org/x/tools/go/packages"
)

type Result struct {
	Contract model.Contract
	Errors   []model.Diagnostic
	Warnings []model.Diagnostic
}

type packageCalls struct {
	pkg   *packages.Package
	calls []sdkCall
}

type reading struct {
	project            Project
	loaded             loadedPackages
	result             Result
	reported           map[model.Diagnostic]bool
	reportedGo         map[*packages.Package]bool
	schemaNames        map[*types.TypeName]string
	schemaDeclarations []schemaDeclaration
	enumNames          map[*types.TypeName]string
	enumDeclarations   map[*types.TypeName]*ast.CallExpr
	unions             []unionField
	schemas            map[*types.TypeName]*structSchema
	derivedSchemas     []model.NamedSchema
	embeddedNames      map[*types.TypeName]bool
	directFields       map[*types.Struct][]*directField
	jsonFields         map[*types.Struct]structFields
	fieldNodes         map[*types.Var]model.Node
}

func ProjectToResult(project Project) Result {
	loaded, diagnostics := projectToLoadedPackages(project)
	if len(diagnostics) > 0 {
		return Result{Errors: diagnostics}
	}
	reading := &reading{
		project:          project,
		loaded:           loaded,
		reported:         map[model.Diagnostic]bool{},
		reportedGo:       map[*packages.Package]bool{},
		schemaNames:      map[*types.TypeName]string{},
		enumNames:        map[*types.TypeName]string{},
		enumDeclarations: map[*types.TypeName]*ast.CallExpr{},
		schemas:          map[*types.TypeName]*structSchema{},
		embeddedNames:    map[*types.TypeName]bool{},
		directFields:     map[*types.Struct][]*directField{},
		jsonFields:       map[*types.Struct]structFields{},
		fieldNodes:       map[*types.Var]model.Node{},
	}
	reading.read()
	return reading.result
}

func (this *reading) read() {
	var declaring []packageCalls
	for _, pkg := range this.loaded.packages {
		if !importsSDK(pkg) {
			continue
		}
		if len(pkg.Errors) > 0 {
			this.reportGoErrors(pkg)
			continue
		}
		declaring = append(declaring, packageCalls{pkg: pkg, calls: packageToCalls(pkg)})
	}
	for _, each := range declaring {
		this.readOptionsInPlace(each)
	}
	this.readCalls(declaring, this.readSchemaDeclaration, "Schema")
	this.readCalls(declaring, this.readEnum, "Enum")
	var project []packageCalls
	for _, each := range declaring {
		if this.loaded.inProject[each.pkg.PkgPath] {
			project = append(project, each)
		}
	}
	this.readProjectStructs()
	this.readCalls(project, this.readResource, "Provides", "Consumes")
	this.readSchemaDeclarationTypes()
	this.readUnionMembers()
	this.schemasToContract()
}

func (this *reading) readCalls(declaring []packageCalls, read func(*packages.Package, sdkCall), functions ...string) {
	for _, each := range declaring {
		for _, call := range each.calls {
			if slices.Contains(functions, call.function) {
				read(each.pkg, call)
			}
		}
	}
}

func (this *reading) location(node ast.Node) model.SourceLocation {
	return this.locationAt(node.Pos())
}

func (this *reading) addError(node ast.Node, format string, arguments ...any) {
	this.addErrorAt(node.Pos(), format, arguments...)
}

func (this *reading) addErrorAt(position token.Pos, format string, arguments ...any) {
	this.appendDiagnostic(&this.result.Errors, model.Diagnostic{Message: fmt.Sprintf(format, arguments...), Location: this.locationAt(position)})
}

func (this *reading) addWarningAt(position token.Pos, format string, arguments ...any) {
	this.appendDiagnostic(&this.result.Warnings, model.Diagnostic{Message: fmt.Sprintf(format, arguments...), Location: this.locationAt(position)})
}

// A diagnostic reached more than once, like a field's from every schema that embeds its struct, is kept once.
func (this *reading) appendDiagnostic(diagnostics *[]model.Diagnostic, diagnostic model.Diagnostic) {
	if this.reported[diagnostic] {
		return
	}
	this.reported[diagnostic] = true
	*diagnostics = append(*diagnostics, diagnostic)
}

func (this *reading) reportGoErrors(pkg *packages.Package) {
	if this.reportedGo[pkg] {
		return
	}
	this.reportedGo[pkg] = true
	for _, diagnostic := range packageToGoErrors(pkg, this.project.Dir) {
		this.appendDiagnostic(&this.result.Errors, diagnostic)
	}
}

func (this *reading) schemaName(name *types.TypeName) string {
	if declared, found := this.schemaNames[name]; found {
		return declared
	}
	return name.Name()
}
