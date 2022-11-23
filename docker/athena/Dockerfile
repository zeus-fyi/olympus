#syntax=docker/dockerfile:1.4
FROM golang:1.19-buster as src_builder

WORKDIR /app

COPY . .

RUN go mod download

COPY configs/ca-certificate.crt .
COPY configs/zeus.fyi/ca.pem .
COPY configs/zeus.fyi/ca.key .

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o athena apps/olympus/athena/main.go

#syntax=docker/dockerfile:1.4
FROM golang:1.19-buster as hypnos_builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o hypnos apps/olympus/hypnos/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=src_builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=src_builder /app/ca.pem /etc/ssl/certs
COPY --from=src_builder /app/ca.key /etc/ssl/certs

COPY --from=src_builder /app/athena /usr/bin/athena
COPY --from=hypnos_builder /app/hypnos /usr/bin/hypnos

EXPOSE 9005