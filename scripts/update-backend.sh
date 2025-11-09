#!/bin/bash

# Cron job that checks for updates and rebuilds if it finds any.

set -euo pipefail

cd "$(dirname "$0")/.."

echo "Checking for updates."
git fetch origin main
LOCAL=$(git rev-parse @)
REMOTE=$(git rev-parse origin/main)

if [ $LOCAL = $REMOTE ]; then
	echo "No changes found."
	exit 0
fi

echo "Changes detected, pulling updates..."
git pull origin main

echo "Rebuilding containers..."
docker compose build --pull # --pull to pull images instead of using cached images.

echo "Restarting services..."
docker compose up -d

echo "Cleaning up old images..."
docker image prune -f

echo "Update completed: $(date)"
