#!/bin/sh

# Set the database URL based on environment variables
export DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE:-disable}"

echo "Running database migrations..."
migrate -path /app/migrations -database "$DB_URL" up

echo "Starting POS Server..."
exec "/app/pos-server"
