#syntax=docker/dockerfile:1.4
FROM golang:1.19-buster as src_builder

WORKDIR /app

COPY --from=src configs/ca-certificate.crt .
COPY --from=src configs/zeus.fyi/ca.pem .
COPY --from=src configs/zeus.fyi/ca.key .

# Download necessary Go modules
COPY --from=src . .

COPY --from=src go.mod .
COPY --from=src go.sum .
COPY --from=src go.work .
COPY --from=src go.work.sum .

RUN go mod download

COPY . .

FROM src_builder AS app_builder

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o athena apps/olympus/athena/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=src_builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=src_builder /app/ca.pem /etc/ssl/certs
COPY --from=src_builder /app/ca.key /etc/ssl/certs

COPY --from=app_builder /app/athena /usr/bin/athena
EXPOSE 9003