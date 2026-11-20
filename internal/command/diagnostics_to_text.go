package command

import (
	"cmp"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/bidirekt/go-sdk/internal/model"
)

type prefixedDiagnostic struct {
	diagnostic model.Diagnostic
	prefix     string
}

func DiagnosticsToText(errors, warnings []model.Diagnostic, cwd string) string {
	prefixed := make([]prefixedDiagnostic, 0, len(errors)+len(warnings))
	for _, diagnostic := range errors {
		prefixed = append(prefixed, prefixedDiagnostic{diagnostic: diagnostic})
	}
	for _, diagnostic := range warnings {
		prefixed = append(prefixed, prefixedDiagnostic{diagnostic: diagnostic, prefix: "warning: "})
	}
	slices.SortStableFunc(prefixed, func(a, b prefixedDiagnostic) int {
		return cmp.Or(
			strings.Compare(a.diagnostic.Location.File, b.diagnostic.Location.File),
			cmp.Compare(a.diagnostic.Location.Line, b.diagnostic.Location.Line),
		)
	})
	var lines []string
	for _, each := range prefixed {
		if line := diagnosticToLine(each.diagnostic, each.prefix, cwd); !slices.Contains(lines, line) {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "")
}

func diagnosticToLine(diagnostic model.Diagnostic, prefix, cwd string) string {
	if diagnostic.Location.File == "" {
		return prefix + diagnostic.Message + "\n"
	}
	return pathToShown(diagnostic.Location.File, cwd) + ":" + strconv.Itoa(diagnostic.Location.Line) + " " + prefix + diagnostic.Message + "\n"
}

func pathToShown(path, cwd string) string {
	relative, err := filepath.Rel(cwd, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(relative)
}
