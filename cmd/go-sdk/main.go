package main

import (
	"fmt"
	"os"

	"github.com/bidirekt/go-sdk/internal/command"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read the current folder: %s\n", err)
		os.Exit(1)
	}
	result := command.ArgsToResult(os.Args[1:], cwd)
	fmt.Fprint(os.Stderr, result.Stderr)
	os.Exit(result.ExitCode)
}
