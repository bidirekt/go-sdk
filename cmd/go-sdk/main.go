package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "Usage: go-sdk generate [--check] [--project <dir>] [--out <file>]")
	os.Exit(1)
}
