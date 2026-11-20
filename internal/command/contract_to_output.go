package command

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bidirekt/go-sdk/internal/model"
	"github.com/bidirekt/go-sdk/internal/writer"
)

func contractToOutput(contract model.Contract, output string, check bool, cwd, diagnostics string) Result {
	yaml := writer.ContractToYaml(contract)
	shown := pathToShown(output, cwd)
	if check {
		return checkOutput(output, shown, yaml, diagnostics)
	}
	if info, err := os.Stat(output); err == nil && info.IsDir() {
		return Result{ExitCode: 1, Stderr: fmt.Sprintf("%scannot write %q: it is a folder\n", diagnostics, shown)}
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o777); err != nil {
		return Result{ExitCode: 1, Stderr: fmt.Sprintf("%scannot write %q: %s\n", diagnostics, shown, errorToReason(err))}
	}
	if err := os.WriteFile(output, []byte(yaml), 0o666); err != nil {
		return Result{ExitCode: 1, Stderr: fmt.Sprintf("%scannot write %q: %s\n", diagnostics, shown, errorToReason(err))}
	}
	return Result{Stderr: diagnostics}
}

func checkOutput(output, shown, yaml, diagnostics string) Result {
	info, err := os.Stat(output)
	if errors.Is(err, fs.ErrNotExist) {
		return Result{ExitCode: 1, Stderr: fmt.Sprintf("%s%q does not exist: run generate without --check\n", diagnostics, shown)}
	}
	if err == nil && info.IsDir() {
		return Result{ExitCode: 1, Stderr: fmt.Sprintf("%scannot read %q: it is a folder\n", diagnostics, shown)}
	}
	written, err := os.ReadFile(output)
	if err != nil {
		return Result{ExitCode: 1, Stderr: fmt.Sprintf("%scannot read %q: %s\n", diagnostics, shown, errorToReason(err))}
	}
	if strings.ReplaceAll(string(written), "\r\n", "\n") != yaml {
		return Result{ExitCode: 1, Stderr: fmt.Sprintf("%s%q is out of date: run generate without --check\n", diagnostics, shown)}
	}
	return Result{Stderr: diagnostics}
}

func errorToReason(err error) string {
	var pathError *fs.PathError
	if errors.As(err, &pathError) {
		return pathError.Err.Error()
	}
	return err.Error()
}
