FROM golang:1.19-buster as snapshot_builder

WORKDIR /app

ARG GOMODCACHE
ARG GOCACHE
ARG GOOS
ARG GOARCH

COPY configs/ca-certificate.crt .
COPY configs/zeus.fyi/ca.pem .
COPY configs/zeus.fyi/ca.key .

# manage dependencies
COPY . .
RUN --mount=type=cache,target=${GOMODCACHE} go mod download
COPY . .
RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o downloader apps/olympus/downloader/main.go

FROM ubuntu:22.04 AS rclone_builder
RUN apt-get update && apt-get install -y ca-certificates curl unzip

WORKDIR /app

RUN curl -O https://downloads.rclone.org/rclone-current-linux-amd64.zip
RUN unzip rclone-current-linux-amd64.zip
RUN cd rclone-*-linux-amd64 && cp rclone /usr/bin/
RUN chown root:root /usr/bin/rclone
RUN chmod 755 /usr/bin/rclone

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=snapshot_builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=snapshot_builder /app/ca.pem /etc/ssl/certs
COPY --from=snapshot_builder /app/ca.key /etc/ssl/certs

COPY --from=snapshot_builder /app/downloader /usr/bin/downloader
COPY --from=rclone_builder /usr/bin/rclone /usr/bin/rclone

RUN mkdir -p ~/.config/rclone
