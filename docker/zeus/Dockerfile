FROM golang:1.18-buster as builder

WORKDIR /app

# Download necessary Go modules

# Download necessary Go modules
COPY apps/autok8s/go.mod .
COPY apps/autok8s/go.sum .
RUN go mod download

COPY .. .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o autok8s apps/autok8s/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/autok8s /usr/bin/autok8s
EXPOSE 9001