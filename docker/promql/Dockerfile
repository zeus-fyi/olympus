#syntax=docker/dockerfile:1.4
FROM golang:1.21-bookworm as builder

WORKDIR /app

ARG GOMODCACHE
ARG GOCACHE
ARG GOOS
ARG GOARCH

COPY go.mod go.sum ./
RUN --mount=type=cache,target=${GOMODCACHE} go mod download

COPY . .
RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o promql apps/olympus/promql/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/promql /usr/bin/promql

EXPOSE 8000