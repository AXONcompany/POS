#!/usr/bin/env bash
set -euo pipefail

# Runs all pending migrations (up).
# Requires: DB_USER, DB_PASSWORD, DB_NAME
# Optional overrides: MIGRATE_HOST (default: db), MIGRATE_PORT (default: 5432), DB_SSLMODE (default: disable)
#
# Note: This runs the migrate container on the same Docker network as your compose `db`
# service, so it can reach Postgres using the service hostname `db`.

: "${DB_USER:?DB_USER is required}"
: "${DB_PASSWORD:?DB_PASSWORD is required}"
: "${DB_NAME:?DB_NAME is required}"

docker compose run --rm migrate up
