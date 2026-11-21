#!/bin/sh
# Runs every scenario of this folder, or the ones named as arguments, as a developer would: go mod download, generate
# and generate --check on the Go of the scenario's .go-version, then bidirekt validate and publish of the generated
# contract on a throwaway broker; a mistake applies its patch to a copy of another scenario first, and a flow scenario
# runs the commands of its steps.txt instead. See README.md.
# COMPOSE_PROJECT_NAME picks the compose project (default go-sdk-scenarios).
set -eu

dir=$(cd "$(dirname "$0")" && pwd)
project=${COMPOSE_PROJECT_NAME:-go-sdk-scenarios}

# A run that mounts a volume keeps its container (no --rm) for the down on exit, which removes it before the volumes:
# with --rm, the daemon removes a container a signal killed only after down tried the volumes, which then stay.
# Compose takes the kept containers of another profile's service for orphans, and its warning would reach the outputs
# the scenarios compare.
compose() {
  COMPOSE_IGNORE_ORPHANS=1 docker compose --progress quiet -p "$project" -f "$dir/compose.yaml" "$@"
}

failed=0
tmp=$(mktemp -d)
# The broker's log tells why a broker step failed; a scenario that fails before the broker doesn't need it.
broker_log=1
# The down and the rm run even when a step before them fails, like the log on an output pipe that closed (| head).
cleanup() {
  set +e
  if [ "$failed" != 0 ] && [ "$broker_log" = 1 ]; then
    compose logs broker >&2
  fi
  # Without every profile, down keeps the volumes that only the services of a profile use.
  compose --profile '*' down -v --remove-orphans
  rm -rf "$tmp"
}
trap 'failed=$?; cleanup' EXIT
# dash and busybox ash skip the EXIT trap when a signal kills the shell, PIPE included when the output pipe closes.
trap 'exit 1' HUP INT PIPE TERM

newline='
'
# Prints the base scenario of the mistake $1, the name of its .patch, or nothing when $1 is not a mistake.
patch_base() {
  for patch in "$dir/$1"/*.patch; do
    if [ -f "$patch" ]; then
      basename "$patch" .patch
    fi
  done
}

# The project scenarios and the mistakes first, then the flows, which deploy to production.
if [ $# = 0 ]; then
  for scenario in "$dir"/*/; do
    scenario=$(basename "$scenario")
    if [ ! -f "$dir/$scenario/steps.txt" ]; then
      set -- "$@" "$scenario"
    fi
  done
  for steps in "$dir"/*/steps.txt; do
    if [ -f "$steps" ]; then
      set -- "$@" "$(basename "$(dirname "$steps")")"
    fi
  done
fi

# The services of the flow scenario $1: its folders with a go.mod.
flow_services() {
  for module in "$dir/$1"/*/go.mod; do
    if [ -f "$module" ]; then
      basename "$(dirname "$module")"
    fi
  done
}

services=$(compose --profile '*' config --services)
# Exits unless the file $1 of testdata/scenarios names a Go version that compose.yaml has a service for.
check_version() {
  if [ ! -f "$dir/$1" ]; then
    echo "testdata/scenarios/$1 does not exist: write in it the Go version to run on, like 1.26" >&2
    exit 1
  fi
  version=$(cat "$dir/$1")
  if ! printf '%s\n' "$services" | grep -qxF "go-$version"; then
    echo "testdata/scenarios/$1 says $version, but testdata/scenarios/compose.yaml has no service go-$version" >&2
    exit 1
  fi
}

for name in "$@"; do
  if [ -f "$dir/$name/steps.txt" ]; then
    flow=$(flow_services "$name")
    if [ -z "$flow" ]; then
      echo "testdata/scenarios/$name/steps.txt has no service to run in: no folder of testdata/scenarios/$name has a go.mod" >&2
      exit 1
    fi
    for service in $flow; do
      check_version "$name/$service/.go-version"
    done
    continue
  fi
  base=$(patch_base "$name")
  if [ -n "$base" ]; then
    if [ "$base" != "${base%%"$newline"*}" ]; then
      echo "testdata/scenarios/$name has more than one .patch: a mistake is one patch on one scenario" >&2
      exit 1
    fi
    if [ ! -f "$dir/$base/go.mod" ] && [ ! -f "$dir/$base/go.work" ]; then
      echo "testdata/scenarios/$name/$base.patch names no scenario: testdata/scenarios/$base has no go.mod nor go.work" >&2
      exit 1
    fi
  elif [ ! -f "$dir/$name/go.mod" ] && [ ! -f "$dir/$name/go.work" ]; then
    echo "no scenario \"$name\": testdata/scenarios/$name has no go.mod, go.work, .patch nor steps.txt" >&2
    exit 1
  fi
  check_version "$name/.go-version"
done

compose up -d --wait

# docker compose run reads all of a stdin that is not a terminal, such as a script piped to sh.
BIDIREKT_API_KEY=$(compose run --rm -T setup </dev/null)
export BIDIREKT_API_KEY
if ! output=$(compose run -T bidirekt create-environment production 2>&1 </dev/null); then
  printf 'bidirekt create-environment failed: %s\n' "$output" >&2
  exit 1
fi
broker_log=0

# Prints why the scenario failed, below its FAIL line.
fail() {
  failures=$((failures + 1))
  echo "FAIL $1 ($(($(date +%s) - started))s)"
  printf '%s\n' "$2" | sed 's/^/  /'
}

# Runs a bidirekt command in the scenario's copy, or prints why it failed and returns 1.
bidirekt() {
  if ! output=$(compose run -T -w "/work/$name" bidirekt "$@" 2>&1 </dev/null); then
    broker_log=1
    fail "$name" "bidirekt $1 failed:
$output"
    return 1
  fi
}

# Runs a command of steps.txt, $2, in the copy of the service $1, and appends what it prints and its exit code to the
# transcript. bidirekt runs in the CLI's image, anything else with sh in the service's Go.
run_step() {
  case $2 in
    bidirekt | 'bidirekt '*) set -- --entrypoint sh -w "/work/$name/$1" bidirekt -c "exec 2>&1; $2" ;;
    *) set -- -w "/sdk/testdata/scenarios/$name/$1" "go-$(cat "$dir/$name/$1/.go-version")" sh -c "exec 2>&1; $2" ;;
  esac
  code=0
  # compose prints an empty line of its own when the command fails: only its messages go to the transcript.
  compose run -T "$@" >>"$tmp/transcript" 2>"$tmp/compose" </dev/null || code=$?
  grep . "$tmp/compose" >>"$tmp/transcript" || true
  if [ -n "$(tail -c 1 "$tmp/transcript")" ]; then
    echo >>"$tmp/transcript"
  fi
  if [ "$code" != 0 ]; then
    echo "[$code]" >>"$tmp/transcript"
  fi
}

# Runs the flow scenario $name: copies it, runs go mod download and go mod tidy in each service, as generate.sh does,
# then each command of steps.txt, and compares the transcript with steps.txt.
run_flow() {
  copy=/sdk/testdata/scenarios/$name
  copied=0
  for service in $(flow_services "$name"); do
    go=go-$(cat "$dir/$name/$service/.go-version")
    if [ "$copied" = 0 ]; then
      # shellcheck disable=SC2016 # $1 and $2 are the arguments of sh -c, expanded in the go service
      if ! output=$(compose run -T "$go" sh -c 'rm -rf "$1" && cp -R "/scenarios/$2" "$1"' sh "$copy" "$name" 2>&1 </dev/null); then
        fail "$name" "copying the scenario failed:
$output"
        return
      fi
      copied=1
    fi
    if ! output=$(compose run -T -w "$copy/$service" "$go" sh -c 'go mod download && go mod tidy' 2>&1 </dev/null); then
      fail "$name" "go mod download and go mod tidy in $service failed:
$output"
      return
    fi
  done

  : >"$tmp/transcript"
  number=0
  while IFS= read -r line || [ -n "$line" ]; do
    number=$((number + 1))
    case $line in
      *'$ '*) service=${line%%'$ '*} ;;
      *) continue ;;
    esac
    case $service in
      '' | *[!a-z0-9_-]*) continue ;;
    esac
    if [ ! -f "$dir/$name/$service/go.mod" ]; then
      fail "$name" "steps.txt:$number: no service \"$service\": testdata/scenarios/$name/$service/go.mod does not exist"
      return
    fi
    printf '%s\n' "$line" >>"$tmp/transcript"
    run_step "$service" "${line#*'$ '}"
  done <"$dir/$name/steps.txt"

  if ! output=$(diff -u -L steps.txt -L transcript "$dir/$name/steps.txt" "$tmp/transcript"); then
    broker_log=1
    fail "$name" "the steps printed another transcript:
$output"
    return
  fi
  passed=$((passed + 1))
  echo "ok $name ($(($(date +%s) - started))s)"
}

passed=0
failures=0
for name in "$@"; do
  started=$(date +%s)
  if [ -f "$dir/$name/steps.txt" ]; then
    run_flow
    continue
  fi

  if ! contract=$(compose run -T "go-$(cat "$dir/$name/.go-version")" sh /scenarios/generate.sh "$name" 2>&1 </dev/null); then
    fail "$name" "$contract"
    continue
  fi

  if [ -n "$contract" ]; then
    participant=$(printf '%s' "$name" | tr - _)
    bidirekt create-participant "$participant" || continue
    bidirekt validate "$contract" || continue
    bidirekt publish "$contract" --participant "$participant" --version 1 || continue
  fi

  passed=$((passed + 1))
  echo "ok $name ($(($(date +%s) - started))s)"
done

echo "$passed passed, $failures failed"
[ "$failures" = 0 ]
