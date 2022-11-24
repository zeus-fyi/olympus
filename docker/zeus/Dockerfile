FROM golang:1.19-buster as builder

WORKDIR /app

# Download necessary Go modules
COPY . .
COPY configs/ca-certificate.crt .
COPY configs/zeus.fyi/ca.pem .
COPY configs/zeus.fyi/ca.key .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o zeus apps/olympus/zeus/main.go

FROM ubuntu:22.04 as doctl
RUN apt-get update && apt-get install -y wget
WORKDIR /app

RUN wget https://github.com/digitalocean/doctl/releases/download/v1.85.0/doctl-1.85.0-linux-amd64.tar.gz -O doctl.tar.gz
RUN tar xf doctl.tar.gz

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=builder /app/ca.pem /etc/ssl/certs
COPY --from=builder /app/ca.key /etc/ssl/certs

COPY --from=builder /app/zeus /usr/bin/zeus
COPY --from=doctl /app/doctl /usr/local/bin
EXPOSE 9001