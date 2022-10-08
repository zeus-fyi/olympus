FROM golang:1.19-buster as builder

WORKDIR /app

COPY configs/DigiCertGlobalRootCA.crt.pem .
# Download necessary Go modules
COPY apps/olympus/apollo/go.mod .
COPY apps/olympus/apollo/go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o beacon_indexer apps/olympus/apollo/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/DigiCertGlobalRootCA.crt.pem /etc/ssl/certs
COPY --from=builder /app/olympus/apollo /usr/bin/apollo
EXPOSE 9000
