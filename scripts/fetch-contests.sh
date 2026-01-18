#!/bin/bash

set -euo pipefail

cd "$(dirname "$0")/.."

docker compose run --build fetcher --remove-orphans

echo "Cleaning up old images..."
docker image prune -f
