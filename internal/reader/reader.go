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
	inProgress         map[*types.TypeName]bool
}

func ProjectToResult(project Project) Result {
	loaded, diagnostics := projectToLoadedPackages(project)
	if len(diagnostics) > 0 {
		return Result{Errors: diagnostics}
	}
	reading := &reading{
		project:          project,
		loaded:           loaded,
		schemaNames:      map[*types.TypeName]string{},
		enumNames:        map[*types.TypeName]string{},
		enumDeclarations: map[*types.TypeName]*ast.CallExpr{},
		schemas:          map[*types.TypeName]*structSchema{},
		embeddedNames:    map[*types.TypeName]bool{},
		directFields:     map[*types.Struct][]*directField{},
		jsonFields:       map[*types.Struct]structFields{},
		fieldNodes:       map[*types.Var]model.Node{},
		inProgress:       map[*types.TypeName]bool{},
	}
	reading.read()
	return reading.result
}

func (this *reading) read() {
	var declaring []packageCalls
	for _, pkg := range this.loaded.packages {
		if _, found := pkg.Imports[sdkPath]; !found {
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
	this.readTagReferences()
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
	this.result.Errors = append(this.result.Errors, model.Diagnostic{Message: fmt.Sprintf(format, arguments...), Location: this.locationAt(position)})
}

func (this *reading) addWarningAt(position token.Pos, format string, arguments ...any) {
	this.result.Warnings = append(this.result.Warnings, model.Diagnostic{Message: fmt.Sprintf(format, arguments...), Location: this.locationAt(position)})
}

func (this *reading) reportGoErrors(pkg *packages.Package) {
	this.result.Errors = append(this.result.Errors, packageToGoErrors(pkg, this.project.Dir)...)
}

func (this *reading) schemaName(name *types.TypeName) string {
	if declared, found := this.schemaNames[name]; found {
		return declared
	}
	return name.Name()
}
