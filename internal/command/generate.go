package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bidirekt/go-sdk/internal/model"
	"github.com/bidirekt/go-sdk/internal/reader"
)

const (
	usage         = "Usage: go-sdk generate [--check] [--project <dir>] [--out <file>]\n"
	defaultOutput = "contracts/generated.yaml"
)

type Result struct {
	ExitCode int
	Stderr   string
}

func ArgsToResult(args []string, cwd string) Result {
	options, valid := argsToOptions(args)
	if !valid {
		return Result{ExitCode: 1, Stderr: usage}
	}
	if real, err := filepath.EvalSymlinks(cwd); err == nil {
		cwd = real
	}
	return generate(options, cwd)
}

func generate(options generateOptions, cwd string) Result {
	dir := cwd
	if options.project != "" {
		dir = relativePathToAbsolute(cwd, options.project)
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			return Result{ExitCode: 1, Stderr: fmt.Sprintf("project %q not found\n", options.project)}
		}
	}
	project, found := reader.DirToProject(dir)
	if !found {
		return Result{ExitCode: 1, Stderr: fmt.Sprintf(
			"no go.mod found in %q or any folder above it: run generate in a folder of a Go module, or pass --project <a folder of the module>\n",
			pathToShown(dir, cwd),
		)}
	}
	read := reader.ProjectToResult(project)
	contract, mergeErrors := model.DeclarationsToContract(read.Contract)
	errors := append(read.Errors, mergeErrors...)
	diagnostics := DiagnosticsToText(errors, read.Warnings, cwd)
	if len(errors) > 0 {
		return Result{ExitCode: 1, Stderr: diagnostics}
	}
	if isEmpty(contract) {
		return Result{ExitCode: 1, Stderr: fmt.Sprintf(
			"%sno github.com/bidirekt/go-sdk declaration found in %q: nothing to generate\n",
			diagnostics, pathToShown(project.Dir, cwd),
		)}
	}
	output := relativePathToAbsolute(project.Dir, defaultOutput)
	if options.out != "" {
		output = relativePathToAbsolute(cwd, options.out)
	}
	return contractToOutput(contract, output, options.check, cwd, diagnostics)
}

func isEmpty(contract model.Contract) bool {
	return len(contract.Provides) == 0 && len(contract.Consumes) == 0 && len(contract.Schemas) == 0
}

func relativePathToAbsolute(base, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(base, path)
}
