#syntax=docker/dockerfile:1.4
FROM golang:1.19-buster as src_builder

WORKDIR /app

COPY . .

COPY configs/web3signer.ephemery/ca-certificate.crt .

ARG GOMODCACHE
ARG GOCACHE

RUN --mount=type=cache,target=${GOMODCACHE} go mod download

COPY . .

RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o aegis apps/olympus/aegis/main.go

FROM ubuntu:22.04 as web3signer
RUN apt-get update && apt-get install -y wget
WORKDIR /app

ARG VERSION

RUN wget https://artifacts.consensys.net/public/web3signer/raw/names/web3signer.tar.gz/versions/${VERSION}/web3signer-${VERSION}.tar.gz
RUN mv web3signer-${VERSION}.tar.gz web3signer
RUN tar xf web3signer

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=src_builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=web3signer /app/web3signer /usr/local/bin

EXPOSE 9000