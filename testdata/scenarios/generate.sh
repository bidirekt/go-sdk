#!/bin/sh
# Runs in a go-<version> service of compose.yaml for the scenario named by $1: copies the project, downloads its
# dependencies, applies the scenario's patch when it is a mistake on another scenario, runs generate and generate
# --check in the copy as a developer would, and compares the results with the scenario's.
# On success prints the path of the generated contract inside the copy, or nothing when the scenario expects exit 1.
# On failure prints why and exits 1.
set -eu
exec 2>&1

name=$1
scenario=/scenarios/$name
copy=/sdk/testdata/scenarios/$name
results=$(mktemp -d)

project=$scenario
patch=
for file in "$scenario"/*.patch; do
  if [ -f "$file" ]; then
    patch=$file
    project=/scenarios/$(basename "$file" .patch)
  fi
done

# The flags of args.txt are split on whitespace, never expanded as globs.
set -f

rm -rf "$copy"
mkdir -p "$copy"
cp -R "$project/." "$copy"
cd "$copy"
rm -f expected.yaml expected-stderr.txt args.txt

args=
if [ -f "$scenario/args.txt" ]; then
  args=$(cat "$scenario/args.txt")
fi
folder=.
# shellcheck disable=SC2086 # args.txt holds flags separated by whitespace
set -- $args
while [ $# -gt 1 ]; do
  if [ "$1" = --project ]; then
    folder=$2
  fi
  shift
done
contract=${folder#./}/contracts/generated.yaml
contract=${contract#./}

# go mod tidy changes nothing on the committed go.mod and go.sum, until the generator's own dependencies move up.
if ! (cd "$folder" && go mod download && go mod tidy) >"$results/mod" 2>&1; then
  echo "go mod download and go mod tidy in $folder failed:"
  cat "$results/mod"
  exit 1
fi

if [ -n "$patch" ]; then
  if ! git apply "$patch" >"$results/patch" 2>&1; then
    echo "applying $name/$(basename "$patch") failed:"
    cat "$results/patch"
    exit 1
  fi
fi

if [ -f "$scenario/expected-stderr.txt" ]; then
  cp "$scenario/expected-stderr.txt" "$results/expected-stderr"
else
  : >"$results/expected-stderr"
fi

files() {
  find . -type f -exec md5sum {} + | sort -k 2
}

# Runs go tool go-sdk generate with the given flags and the scenario's, and fails unless stdout is empty and stderr
# is expected. Leaves the exit code in $code.
generate() {
  command="go tool go-sdk generate${*:+ $*}${args:+ $args}"
  code=0
  # shellcheck disable=SC2086 # args.txt holds flags separated by whitespace
  go tool go-sdk generate "$@" $args >"$results/stdout" 2>"$results/stderr" </dev/null || code=$?
  if [ -s "$results/stdout" ]; then
    echo "$command printed to stdout:"
    cat "$results/stdout"
    exit 1
  fi
  if ! diff -u -L expected-stderr.txt -L stderr "$results/expected-stderr" "$results/stderr"; then
    echo "$command printed another stderr (exit $code)"
    exit 1
  fi
}

files >"$results/before"

if [ ! -f "$scenario/expected.yaml" ]; then
  generate
  if [ "$code" != 1 ]; then
    echo "$command exited with $code, expected 1"
    exit 1
  fi
  files >"$results/after"
  if ! diff -u -L before -L after "$results/before" "$results/after"; then
    echo "$command wrote files, expected none"
    exit 1
  fi
  exit 0
fi

generate
if [ "$code" != 0 ]; then
  echo "$command exited with $code, expected 0"
  exit 1
fi
if [ ! -f "$contract" ]; then
  echo "$command wrote no $contract"
  exit 1
fi
if ! diff -u -L expected.yaml -L "$contract" "$scenario/expected.yaml" "$contract"; then
  echo "$contract differs from expected.yaml"
  exit 1
fi
files | grep -v "  ./$contract\$" >"$results/after" || true
grep -v "  ./$contract\$" "$results/before" >"$results/before-but-contract" || true
if ! diff -u -L before -L after "$results/before-but-contract" "$results/after"; then
  echo "$command wrote files other than $contract"
  exit 1
fi

generate --check
if [ "$code" != 0 ]; then
  echo "$command exited with $code, expected 0"
  exit 1
fi

echo "$contract"
