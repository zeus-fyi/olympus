#syntax=docker/dockerfile:1.4
FROM python:alpine3.14
WORKDIR /app

RUN apk add --update gcc libc-dev linux-headers

RUN git clone https://github.com/ethereum/staking-deposit-cli.git --branch v2.3.0 --single-branch

COPY requirements.txt setup.py ./

COPY staking_deposit ./staking_deposit

RUN pip3 install -r requirements.txt
RUN python3 setup.py install

ARG cli_command

ENTRYPOINT [ "python3", "./staking_deposit/deposit.py" ]


FROM golang:1.19-buster as src_builder

WORKDIR /app

COPY . .
COPY configs/ca-certificate.crt .
COPY configs/zeus.fyi/ca.pem .
COPY configs/zeus.fyi/ca.key .

ARG GOMODCACHE
ARG GOCACHE

RUN --mount=type=cache,target=${GOMODCACHE} go mod download

COPY . .

RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o hera apps/olympus/hera/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=src_builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=src_builder /app/ca.pem /etc/ssl/certs
COPY --from=src_builder /app/ca.key /etc/ssl/certs

COPY --from=src_builder /app/hera /usr/bin/hera
EXPOSE 9004