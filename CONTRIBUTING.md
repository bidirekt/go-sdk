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

To run the tests by hand, run `go test ./...` in each module:

```sh
go test ./...
(cd cmd/go-sdk && go test ./...)
```

## Releases

A release is two tags on the same commit: `vX.Y.Z` for the root module and `cmd/go-sdk/vX.Y.Z` for the generator, with the same version. Push both at once:

```sh
git tag vX.Y.Z
git tag cmd/go-sdk/vX.Y.Z
git push --atomic origin vX.Y.Z cmd/go-sdk/vX.Y.Z
```

The `Release` workflow runs on `vX.Y.Z`: it checks that `cmd/go-sdk/vX.Y.Z` points at the same commit, creates the GitHub Release, and requests both versions from proxy.golang.org so they show on pkg.go.dev.

## Security

Report a vulnerability as [SECURITY.md](SECURITY.md) describes, never in an issue.
