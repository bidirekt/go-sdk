#!/bin/sh
# Runs in the go service of compose.yaml, in the folder of a module: downloads its dependencies and installs the
# generator with the commands of the README's Installation, as a developer would. On a module that already has the
# tool, they only bring go.mod and go.sum up to the generator of the checkout.
set -eu

go mod download
go get github.com/bidirekt/go-sdk
go get -tool github.com/bidirekt/go-sdk/cmd/go-sdk
go mod tidy
