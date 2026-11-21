# Bidirekt Go SDK

[![CI](https://img.shields.io/github/actions/workflow/status/bidirekt/go-sdk/ci.yml?branch=main&label=CI)](https://github.com/bidirekt/go-sdk/actions/workflows/ci.yml) [![License](https://img.shields.io/github/license/bidirekt/go-sdk)](LICENSE)

In Bidirekt, each side of an integration declares in a contract file what it provides and what it consumes. This SDK writes that file from the code of a service on Go 1.26 or later: you declare in Go code what the service provides and consumes, mark the fields of its schemas with the `bidirekt` struct tag, run `go tool go-sdk generate`, and get `contracts/generated.yaml`, ready for `bidirekt validate` and `bidirekt publish`. It is the service's only contract file: once a service adopts the SDK, every resource and schema it provides or consumes is declared in its code.

## Development

Read [CONTRIBUTING.md](CONTRIBUTING.md). From the repository root, turn on the `pre-push` hook, which runs the Go checks of CI before every push: `gofmt` from the toolchain in `go.mod`, `go vet`, `golangci-lint`, `go-arch-lint`, `go test` and `govulncheck`.

```sh
git config core.hooksPath .githooks
```

## License

MIT, see [LICENSE](LICENSE).
