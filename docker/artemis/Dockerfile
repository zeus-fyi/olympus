#syntax=docker/dockerfile:1.4
FROM golang:1.19-buster as src_builder

WORKDIR /app

COPY . .
COPY configs/ca-certificate.crt .
COPY configs/zeus.fyi/ca.pem .
COPY configs/zeus.fyi/ca.key .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o artemis apps/olympus/artemis/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=src_builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=src_builder /app/ca.pem /etc/ssl/certs
COPY --from=src_builder /app/ca.key /etc/ssl/certs

COPY --from=src_builder /app/artemis /usr/bin/artemis
EXPOSE 9004