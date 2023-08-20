FROM golang:1.21-bookworm as builder

WORKDIR /app

COPY configs/apollo.zeus.fyi/ca-certificate.crt .
# Download necessary Go modules
COPY apps/olympus/apollo/go.mod .
COPY apps/olympus/apollo/go.sum .

ARG GOMODCACHE
ARG GOCACHE

RUN --mount=type=cache,target=${GOMODCACHE} go mod download
COPY . .

RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o apollo apps/olympus/apollo/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=builder /app/apollo /usr/bin/apollo
EXPOSE 9000
