#syntax=docker/dockerfile:1.4
FROM golang:1.21-bookworm as src_builder

WORKDIR /app

COPY . .

ARG GOMODCACHE
ARG GOCACHE
ARG GOOS
ARG GOARCH

RUN --mount=type=cache,target=${GOMODCACHE} go mod download

COPY . .

RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o hephaestus apps/olympus/hephaestus/main.go

