# Fixtures

Each folder is one case of `generate`, run by `tests/fixtures/fixtures_test.go`. Adding a case is adding a folder.

The runner copies the folder to a temporary directory, points the `replace` of `github.com/bidirekt/go-sdk` at this checkout, and runs `generate` there with `GOWORK=off`, that directory as the current folder, so nothing is written into the repository.

A case holds:

- A Go module: a `go.mod` with `require github.com/bidirekt/go-sdk v0.0.0` and `replace github.com/bidirekt/go-sdk => ../../../../..`, and the code. A module that only imports the SDK and the standard library needs no `go.sum`.
- The expected result, one of:
  - `expected.yaml`: `generate` exits 0, the contract equals this file byte by byte, and stderr is empty.
  - `expected.yaml` and `expected-stderr.txt`: the same, with stderr equal to `expected-stderr.txt` (warnings).
  - `expected-stderr.txt` alone: `generate` exits 1, stderr equals this file byte by byte, and no file of the case changes.
- `args.txt` (optional): the flags passed after `generate`, separated by whitespace, e.g. `--check`. The contract is read from `contracts/generated.yaml` of the project folder: the case folder, or the folder `--project` names.

For `--check`, commit the `contracts/generated.yaml` the check reads inside the case: up to date pairs it with an identical `expected.yaml`, out of date with an `expected-stderr.txt`.
