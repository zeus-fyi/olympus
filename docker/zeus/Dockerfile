FROM golang:1.19-buster as builder

WORKDIR /app

COPY configs/ca-certificate.crt .
# Download necessary Go modules
COPY apps/olympus/zeus/go.mod .
COPY apps/olympus/zeus/go.sum .
RUN go mod download

COPY .. .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o zeus apps/olympus/zeus/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=builder /app/zeus /usr/bin/zeus
EXPOSE 9001