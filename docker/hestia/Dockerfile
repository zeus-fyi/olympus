#syntax=docker/dockerfile:1.4
FROM golang:1.19-buster as src_builder

WORKDIR /app

COPY . .
COPY configs/ca-certificate.crt .
COPY configs/zeus.fyi/ca.pem .
COPY configs/zeus.fyi/ca.key .

ARG GOMODCACHE
ARG GOCACHE

RUN --mount=type=cache,target=${GOMODCACHE} go mod download

COPY . .

RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o hestia apps/olympus/hestia/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=src_builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=src_builder /app/ca.pem /etc/ssl/certs
COPY --from=src_builder /app/ca.key /etc/ssl/certs

COPY --from=src_builder /app/hestia /usr/bin/hestia
EXPOSE 9002
