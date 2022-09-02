#!/bin/sh

set -e

cd "$(dirname "$0")/.."

[ -z "$DEBUG" ] || set -x

echo "==> Running example-app tests…"
if [ -n "$1" ]; then
  testdir="/$*"
fi

docker-compose run --rm -e EXPAPP_PORT=8080 web go test -cover -race .$testdir/...
