FROM golang:1.19-buster as builder

WORKDIR /app

# Download necessary Go modules
COPY . .

RUN go mod download

COPY configs/ca-certificate.crt .
COPY configs/zeus.fyi/ca.pem .
COPY configs/zeus.fyi/ca.key .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o athena apps/olympus/athena/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=builder /app/ca.pem /etc/ssl/certs
COPY --from=builder /app/ca.key /etc/ssl/certs

COPY --from=builder /app/athena /usr/bin/athena
EXPOSE 9003