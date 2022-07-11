FROM golang:1.18-buster as builder

WORKDIR /app

# Download necessary Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o olympus

FROM ubuntu:22.04
COPY --from=builder /app/olympus /app/
EXPOSE 8080
ENTRYPOINT ["/app/olympus"]