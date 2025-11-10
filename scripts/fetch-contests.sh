#!/bin/bash

docker compose run --build fetcher

echo "Cleaning up old images..."
docker image prune -f
