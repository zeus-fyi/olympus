#syntax=docker/dockerfile:1.4
FROM ubuntu:22.04 as t-digest-builder
WORKDIR /app

ARG VERSION
RUN apt-get update && apt-get install -y ca-certificates git
RUN git clone https://github.com/usmanm/redis-tdigest.git

WORKDIR /app/redis-tdigest
RUN git pull origin master
RUN make all
