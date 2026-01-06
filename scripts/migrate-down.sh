#!/usr/bin/env bash
set -euo pipefail

# Rolls back N migrations (default: 1).
# Requires: DB_USER, DB_PASSWORD, DB_NAME
# Optional overrides: MIGRATE_HOST (default: db), MIGRATE_PORT (default: 5432), DB_SSLMODE (default: disable)
#
# Note: This runs the migrate container on the same Docker network as your compose `db`
# service, so it can reach Postgres using the service hostname `db`.

: "${DB_USER:?DB_USER is required}"
: "${DB_PASSWORD:?DB_PASSWORD is required}"
: "${DB_NAME:?DB_NAME is required}"

STEPS="${1:-1}"

docker compose run --rm migrate down "$STEPS"
