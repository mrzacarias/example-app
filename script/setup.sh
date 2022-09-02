#!/bin/sh

set -e

cd "$(dirname "$0")/.."

echo "==> Building docker container and initializing everything..."
docker-compose build

script/update.sh

echo "==> App is now ready to go!"