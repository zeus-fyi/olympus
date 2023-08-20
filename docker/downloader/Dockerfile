#syntax=docker/dockerfile:1.4
FROM golang:1.21-bookworm as snapshot_builder

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

COPY . .
RUN --mount=type=cache,target=${GOMODCACHE} go mod download
COPY . .
RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o downloader apps/olympus/downloader/main.go

FROM ubuntu:22.04 AS rclone_builder
RUN apt-get update && apt-get install -y ca-certificates curl unzip

WORKDIR /app

RUN curl -O https://downloads.rclone.org/rclone-current-linux-amd64.zip
RUN unzip rclone-current-linux-amd64.zip
RUN cd rclone-*-linux-amd64 && cp rclone /usr/bin/
RUN chown root:root /usr/bin/rclone
RUN chmod 755 /usr/bin/rclone

FROM ethereum/client-go:v1.11.6 as geth
WORKDIR /app

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=snapshot_builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=snapshot_builder /app/ca.pem /etc/ssl/certs
COPY --from=snapshot_builder /app/ca.key /etc/ssl/certs

COPY --from=snapshot_builder /app/downloader /usr/bin/downloader
COPY --from=rclone_builder /usr/bin/rclone /usr/bin/rclone
COPY --from=geth /usr/local/bin/geth /usr/local/bin

RUN mkdir -p ~/.config/rclone
