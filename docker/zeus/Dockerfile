#syntax=docker/dockerfile:1.4
FROM golang:1.20.1-buster as builder

WORKDIR /app
# Download necessary Go modules
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
ARG GOOS
ARG GOARCH

RUN --mount=type=cache,target=${GOMODCACHE} go mod download

COPY . .

RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o zeus ~/apps/olympus/zeus/main.go

FROM ubuntu:22.04 as doctl
RUN apt-get update && apt-get install -y wget
WORKDIR /app

RUN wget https://github.com/digitalocean/doctl/releases/download/v1.85.0/doctl-1.85.0-linux-amd64.tar.gz -O doctl.tar.gz
RUN tar xf doctl.tar.gz

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=builder /app/ca.pem /etc/ssl/certs
COPY --from=builder /app/ca.key /etc/ssl/certs

COPY --from=builder /app/zeus /usr/bin/zeus
COPY --from=doctl /app/doctl /usr/local/bin
EXPOSE 9001