# Contributing to the Bidirekt Go SDK

## Bugs and feature requests

Open an issue in this repository, with the **Bug report** or the **Feature request** form.

## Questions

For questions about using Bidirekt, read the docs at https://bidirekt.com/docs.

## Pull requests

Fork the repository, create a branch, and open a pull request against `main`. The CI checks `fmt-vet`, `lint`, `test`, `govulncheck`, `arch-lint`, `shellcheck`, `e2e` and `scenarios` must pass before it is merged. The repository holds two Go modules: the root, `github.com/bidirekt/go-sdk` (the `bidirekt` package), and `cmd/go-sdk`, `github.com/bidirekt/go-sdk/cmd/go-sdk` (the generator). From the repository root, turn on the `pre-push` hook, which runs `gofmt` once and then, in each module, `go vet`, `golangci-lint`, `go-arch-lint`, `go test` and `govulncheck` before every push:

```sh
git config core.hooksPath .githooks
```

## Security

Report a vulnerability as [SECURITY.md](SECURITY.md) describes, never in an issue.
