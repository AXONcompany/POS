#!/usr/bin/env bash
set -euo pipefail

# Usage:
#   source scripts/load-db-env.sh            # loads from .env if present, otherwise prompts
#   source scripts/load-db-env.sh .env.local # loads from the given env file

ENV_FILE="${1:-.env}"

if [[ -f "$ENV_FILE" ]]; then
  set -a
  source "$ENV_FILE"
  set +a
fi

: "${DB_NAME:=${DB_NAME:-}}"
: "${DB_USER:=${DB_USER:-}}"
: "${DB_PASSWORD:=${DB_PASSWORD:-}}"

if [[ -z "${DB_NAME}" ]]; then
  read -r -p "DB_NAME: " DB_NAME
fi

if [[ -z "${DB_USER}" ]]; then
  read -r -p "DB_USER: " DB_USER
fi

if [[ -z "${DB_PASSWORD}" ]]; then
  read -r -s -p "DB_PASSWORD: " DB_PASSWORD
  echo
fi

export DB_NAME DB_USER DB_PASSWORD

echo "Exported DB_NAME and DB_USER (DB_PASSWORD hidden)."
