#!/bin/sh
# Publishes the expected.yaml of every success fixture to a throwaway broker, one participant per fixture.
# Proves the broker accepts what generate writes. COMPOSE_PROJECT_NAME picks the compose project (default go-sdk-e2e).
set -eu

dir=$(cd "$(dirname "$0")" && pwd)
fixtures=$(cd "$dir/../../cmd/go-sdk/testdata/fixtures" && pwd)
project=${COMPOSE_PROJECT_NAME:-go-sdk-e2e}

compose() {
  docker compose -p "$project" -f "$dir/compose.yaml" "$@"
}

failed=0
# The down runs even when a step before it fails, like the log on an output pipe that closed (| head).
cleanup() {
  set +e
  if [ "$failed" != 0 ]; then
    compose logs broker >&2
  fi
  compose down -v --remove-orphans
}
trap 'failed=$?; cleanup' EXIT
# dash and busybox ash skip the EXIT trap when a signal kills the shell, PIPE included when the output pipe closes.
trap 'exit 1' HUP INT PIPE TERM

compose up -d --wait

# docker compose run reads all of a stdin that is not a terminal, such as a script piped to sh.
BIDIREKT_API_KEY=$(compose run --rm -T setup </dev/null)
export BIDIREKT_API_KEY

published=0
rejected=0
for expected in "$fixtures"/*/expected.yaml; do
  fixture=$(basename "$(dirname "$expected")")
  participant=e2e_$(printf '%s' "$fixture" | tr - _)

  if output=$(compose run --rm -T cli create-participant "$participant" 2>&1 </dev/null) &&
    output=$(compose run --rm -T cli publish "$fixture/expected.yaml" --participant "$participant" --version 1 2>&1 </dev/null); then
    published=$((published + 1))
    echo "ok $fixture"
  else
    rejected=$((rejected + 1))
    echo "FAIL $fixture"
    printf '%s\n' "$output" | sed 's/^/  /'
  fi
done

echo "$published published, $rejected rejected"
[ "$rejected" = 0 ]
