#syntax=docker/dockerfile:1.4
FROM golang:1.21-bookworm as builder

WORKDIR /app

ARG GOMODCACHE
ARG GOCACHE
ARG GOOS
ARG GOARCH

ARG CA_CRT
RUN if [ -z "${CA_CRT}" ]; then \
      COPY configs/ca-certificate.crt .; \
    else \
      echo "${CA_CRT}" > ca-certificate.crt; \
    fi

ARG CA_PEM
RUN if [ -z "${CA_PEM}" ]; then \
      COPY configs/zeus.fyi/ca.pem .; \
    else \
      echo "${CA_PEM}" > ca.pem; \
    fi

ARG CA_KEY
RUN if [ -z "${CA_KEY}"]; then \
      COPY configs/zeus.fyi/ca.key .; \
    else \
      echo "${CA_KEY}" > ca.key; \
    fi

COPY go.mod go.sum ./
RUN --mount=type=cache,target=${GOMODCACHE} go mod download

COPY . .
RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o hypnos apps/olympus/hypnos/main.go

#syntax=docker/dockerfile:1.4
FROM ghcr.io/foundry-rs/foundry:latest as foundry

WORKDIR /app

COPY /usr/local/bin/anvil .

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=builder /app/ca.pem /etc/ssl/certs
COPY --from=builder /app/ca.key /etc/ssl/certs

COPY --from=builder /app/hypnos /usr/bin/hypnos
COPY --from=foundry /app/anvil /usr/bin/anvil

EXPOSE 8888, 8545