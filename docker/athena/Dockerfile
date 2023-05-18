#syntax=docker/dockerfile:1.4
FROM golang:1.20.1-buster as src_builder

WORKDIR /app

COPY . .

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

RUN --mount=type=cache,target=${GOMODCACHE} go mod download

COPY . .

RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o athena apps/olympus/athena/main.go

# Fetch the source code
FROM golang:1.20-alpine as clone
RUN apk add --no-cache git
ARG VERSION
RUN git clone https://github.com/ethereum/go-ethereum.git --branch ${VERSION} --single-branch /go-ethereum

# Build Geth in a stock Go builder container
FROM golang:1.20-alpine as p2p
RUN apk add --no-cache gcc musl-dev linux-headers git

# Get dependencies - will also be cached if we won't change go.mod/go.sum
COPY --from=clone /go-ethereum/go.mod /go-ethereum/
COPY --from=clone /go-ethereum/go.sum /go-ethereum/
RUN cd /go-ethereum && go mod download

ADD . /go-ethereum

COPY --from=clone /go-ethereum/build /go-ethereum/
RUN cd /go-ethereum && go run build/ci.go install -static ./cmd/devp2p

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=src_builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=src_builder /app/ca.pem /etc/ssl/certs
COPY --from=src_builder /app/ca.key /etc/ssl/certs

COPY --from=src_builder /app/athena /usr/bin/athena
COPY --from=p2p /go-ethereum/build/bin/devp2p /usr/local/bin/

EXPOSE 9003