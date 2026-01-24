#!/bin/bash

set -euo pipefail

cd "$(dirname "$0")/.."

docker compose run --build --rm --remove-orphans fetcher "$@" # "$@" forwards all arguments.

echo "Cleaning up old images..."
docker image prune -f
