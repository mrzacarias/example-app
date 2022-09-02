#!/bin/sh

set -e

cd "$(dirname "$0")/.."

echo "==> Deleting all containers, volumes, networks and images..."
docker-compose -f ./docker-compose.yml down --rmi=local --volumes --remove-orphans
docker system prune --volumes -f

echo "==> App was nuked!"
