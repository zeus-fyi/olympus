FROM golang:1.19-buster as builder

WORKDIR /app

ARG GITHUB_TOKEN
ENV TOKEN=$GITHUB_TOKEN

RUN go env -w GOPRIVATE="github.com/zeus-fyi/olympus,github.com/zeus-fyi/olympus/zeus"
RUN git config --global url."https://${TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN rm go.mod
RUN rm go.sum

COPY configs/ca-certificate.crt .
COPY configs/zeus.fyi/ca.pem .
COPY configs/zeus.fyi/ca.key .

# Download necessary Go modules
COPY apps/olympus/zeus/go.mod .
COPY apps/olympus/zeus/go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o zeus apps/olympus/zeus/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=builder /app/ca.pem /etc/ssl/certs
COPY --from=builder /app/ca.key /etc/ssl/certs

COPY --from=builder /app/zeus /usr/bin/zeus
EXPOSE 9001