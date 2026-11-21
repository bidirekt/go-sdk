# Scenarios

Each folder is a project as a developer has it, run end to end by `run.sh`: the generator installed as the README's Installation tells, `go tool go-sdk generate`, `generate --check`, then `bidirekt validate` and `bidirekt publish` of the generated contract on the broker pinned in `../e2e/compose.yaml`. What the scenario checks is what the developer sees, so the SDK may change in any way inside as long as every scenario passes. Adding a scenario is adding a folder.

A mistake is a patch on a copy of a project scenario: the mistake a developer makes in that project, and what `generate` tells them.

A flow scenario is two teams living with each other's contracts: several projects, a provider and its consumers, and the commands their developers type, in order, with what each one prints, from `generate` to the `can-i-deploy` that fails when the provider breaks a field a consumer reads.

Run them all, or some by name, with only Docker installed:

```sh
sh testdata/scenarios/run.sh
sh testdata/scenarios/run.sh chi
```

Go runs in the pinned `golang` image of the `go` service of `compose.yaml`, never on the machine: a run on a laptop and the `scenarios` job of CI use the same Go, and a new Go version is a new image there.

## A scenario

- The project: a Go module, `go.mod` and its committed `go.sum`, or a workspace, `go.work` and the modules it uses. The frameworks come from the Go proxy; the SDK comes from the checkout, through two `replace` directives in the `go.mod` of each module that uses it: `github.com/bidirekt/go-sdk => ../../..` and `github.com/bidirekt/go-sdk/cmd/go-sdk => ../../../cmd/go-sdk`, a folder deeper for each folder the module is below the scenario's. The module that runs `generate` has the generator as a tool, installed as below.
- `args.txt` (optional): the flags after `generate`, separated by whitespace, like `--project apps/api`. The contract is read from `contracts/generated.yaml` of the `--project` folder, or of the scenario's.
- What the developer must see:
  - `expected.yaml`: `generate` exits 0 and writes only the contract, equal to this file byte by byte.
  - No `expected.yaml`: `generate` exits 1 and writes nothing.
  - `expected-stderr.txt`: the stderr of `generate`, byte by byte, with or without `expected.yaml` (warnings, or the errors); without the file, stderr is empty. stdout is always empty.

Name the folder in kebab-case: its participant on the broker is the name in snake_case.

## A mistake

- `<scenario>.patch`: the change, in the format of `git diff`, applied with `git apply` to a copy of the project scenario `<scenario>`, like `net-http.patch`. One patch per mistake.
- `args.txt` and `expected-stderr.txt` as in a scenario, and no `expected.yaml`: a mistake stops `generate`, which must say where and how to fix it.

Name the folder `mistake-<what the developer did>`. To write a patch, change a copy of the project in a folder with git and run `git diff`.

## A flow scenario

- A folder per service, each a Go module as above, with its `go.mod` and `go.sum`; its `replace` directives point a folder deeper, `../../../..`. `run.sh` installs the generator in each one before the steps.
- `steps.txt`: the terminal of the developers, as they would see it. A line `<service>$ <command>` is a command typed in the folder of that service; the lines below it, up to the next command, are what it prints, stdout and stderr together; `[<code>]` closes the output of a command that exits with a code other than 0. A command that starts with `bidirekt` runs in the CLI pinned in the `bidirekt` service of `compose.yaml`, with the broker of the run; any other runs with `sh` in the `go` service, where `go`, `cat` and `git` are. Nothing else in the file: no comments, no blank lines but the ones a command prints.
- Anything else the steps use, like a patch for `git apply ../<name>.patch`, beside `steps.txt`.

`cross-service` is the example: `petstore-api` (`net/http`) publishes v1 and records its deployment to `production`; `storefront` (an `http.Client` with `bidirekt.Consumes` next to each call) validates, publishes, gets a green `can-i-deploy` against the deployed `petstore-api`, and records its deployment; then `rename-pet-name.patch` renames the `name` that `storefront` reads, and the regenerated `petstore-api` gets a red `validate` and `can-i-deploy` with the break.

A flow of one service shows what the broker answers to a contract `generate` writes: `mistake-path-parameter` declares `/adoptions/{petId}`, the pattern of its `http.ServeMux` route, which `generate` writes as is, and `validate` refuses.

To write the steps, write the commands alone and run the scenario: the diff shows what each one printed, to read and copy in.

The participants are the names the steps give them, and every scenario of a run shares the broker: give each flow participants no other scenario uses.

The project files are code of an example project, never built by the tests: `gofmt`, `go vet`, golangci-lint, go-arch-lint and govulncheck skip `testdata/`.

## What `run.sh` does

It starts the broker of `../e2e/compose.yaml`, creates the API key and the `production` environment, and then, per project scenario and mistake:

1. In the `go` service, copies the project, without `expected.yaml`, `expected-stderr.txt` and `args.txt`, to the same path over the checkout (so `../../..` is the SDK), and installs the generator in the `--project` folder, or the scenario's, with the commands of the README's Installation: `go mod download`, `go get github.com/bidirekt/go-sdk`, `go get -tool github.com/bidirekt/go-sdk/cmd/go-sdk` and `go mod tidy`. On the committed `go.mod` and `go.sum` they change nothing, until the generator's own dependencies move up. For a mistake, it then applies the patch to the copy. The scenarios of a run share the module cache and the build cache.
2. Runs `generate` with `args.txt` and compares exit code, stdout, stderr, the contract and the files written; with `expected.yaml`, runs `generate --check`, which must exit 0 with the same stderr.
3. With `expected.yaml`: `bidirekt create-participant`, `bidirekt validate` and `bidirekt publish --version 1` of the generated contract.

Then, per flow scenario (run without names, the flows come after the project scenarios and the mistakes), it copies the folder in the same way, installs the generator in each service, runs each command of `steps.txt` in the copy of its service, and compares what they printed with `steps.txt`.

It prints `ok <scenario>` or `FAIL <scenario>` with the reason (the diff when the contract, stderr or transcript differs), then a summary line, and exits 1 when a scenario fails. It prints the broker's log only when a broker step failed or a flow's transcript differs. On exit and on a signal it removes the compose project with its volumes.

A project scenario's `validate` runs without `--participant` and `--environment`: it checks only the file, so a project that consumes passes with nothing deployed. To compare a consumer with its deployed provider, write a flow scenario whose steps publish the provider and record its deployment first.

## Writing go.mod and go.sum

Write the project's `go.mod` with the `module` line, the `go` line and the two `replace` directives, then, from the repository root, run the Installation commands in the `go` service, with the project's folder mounted over the checkout:

```sh
docker compose -p go-sdk-scenarios-mod -f testdata/scenarios/compose.yaml run --rm -T \
  -v "$PWD/testdata/scenarios/chi:/sdk/testdata/scenarios/chi" -w /sdk/testdata/scenarios/chi \
  go sh -c 'go get github.com/go-chi/chi/v5 && sh /scenarios/install.sh' </dev/null
docker compose -p go-sdk-scenarios-mod -f testdata/scenarios/compose.yaml --profile '*' down -v
```

For a module below the scenario's folder, like a service of a flow scenario or a module of a workspace, mount the scenario's folder and work in the module's.
