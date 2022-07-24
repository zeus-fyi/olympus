FROM golang:1.18-buster as builder

WORKDIR /app

COPY configs/ca-certificate.crt .
# Download necessary Go modules
COPY apps/beacon_indexer/go.mod .
COPY apps/beacon_indexer/go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o beacon_indexer apps/beacon_indexer/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=builder /app/beacon_indexer /usr/bin/beacon_indexer
EXPOSE 9000
