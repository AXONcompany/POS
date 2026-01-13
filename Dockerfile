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

RUN adduser -D -H -s /sbin/nologin appuser
WORKDIR /app

COPY --from=build /out/pos-server /app/pos-server

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/pos-server"]
