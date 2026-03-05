# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build

WORKDIR /src

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/pos-server ./cmd/server


FROM alpine:3.20

RUN apk add --no-cache curl tar \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate \
    && adduser -D -H -s /sbin/nologin appuser

WORKDIR /app

COPY --from=build /out/pos-server /app/pos-server
COPY db/migrations /app/migrations
COPY docker-entrypoint.sh /app/docker-entrypoint.sh

RUN chmod +x /app/docker-entrypoint.sh \
    && chown -R appuser:appuser /app

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/docker-entrypoint.sh"]
