#syntax=docker/dockerfile:1.4
FROM golang:1.21-bookworm as src_builder

WORKDIR /app

COPY . .
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

ARG GOMODCACHE
ARG GOCACHE

COPY go.mod go.sum ./

RUN --mount=type=cache,target=${GOMODCACHE} go mod download

COPY . .

RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o hera apps/olympus/hera/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=src_builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=src_builder /app/ca.pem /etc/ssl/certs
COPY --from=src_builder /app/ca.key /etc/ssl/certs

COPY --from=src_builder /app/hera /usr/bin/hera
EXPOSE 9008