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

# Build Devp2p in a stock Go builder container
FROM golang:1.20-alpine as p2p

WORKDIR /app
ARG VERSION
RUN apk add --no-cache gcc musl-dev linux-headers git

# Set the current working directory inside the container.
WORKDIR /app

ARG VERSION

# Clone the specific version of the go-ethereum repository.
RUN git clone --depth 1 --branch ${VERSION} https://github.com/ethereum/go-ethereum.git

# Copy the files from the host to the current location (i.e., '/app') in the container.
COPY ./go-ethereum/go.mod ./go-ethereum/go.sum ./go-ethereum/

# Download go module dependencies.
RUN cd go-ethereum && go mod download

# Copy all the files from the host to the current location (i.e., '/app/go-ethereum') in the container.
COPY . ./go-ethereum

# Build the application.
RUN cd go-ethereum && go run build/ci.go install -static ./cmd/devp2p

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=src_builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=src_builder /app/ca.pem /etc/ssl/certs
COPY --from=src_builder /app/ca.key /etc/ssl/certs

COPY --from=src_builder /app/athena /usr/bin/athena
COPY --from=p2p /go-ethereum/build/bin/devp2p /usr/local/bin/

EXPOSE 9003