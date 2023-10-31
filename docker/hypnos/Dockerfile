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
RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o hypnos apps/olympus/hypnos/main.go

#syntax=docker/dockerfile:1.4
FROM ghcr.io/foundry-rs/foundry:latest as foundry

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=builder /app/ca.pem /etc/ssl/certs
COPY --from=builder /app/ca.key /etc/ssl/certs

COPY --from=builder /app/hypnos /usr/bin/hypnos

COPY --from=foundry /usr/local/bin/anvil /usr/bin/anvil

EXPOSE 8888 8545