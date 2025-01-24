#syntax=docker/dockerfile:1.4
FROM golang:1.21-bookworm as builder

WORKDIR /app

ARG CA_CRT
RUN if [ -z "${CA_CRT}" ]; then \
      COPY configs/ca-certificate.crt .; \
    else \
      echo "${CA_CRT}" > ca-certificate.crt; \
    fi

ARG CA_PEM
RUN if [ -z "${CA_PEM}" ]; then \
      COPY configs/zeus.fyi/ca.pem .; \
    else \
      echo "${CA_PEM}" > ca.pem; \
    fi

ARG CA_KEY
RUN if [ -z "${CA_KEY}"]; then \
      COPY configs/zeus.fyi/ca.key .; \
    else \
      echo "${CA_KEY}" > ca.key; \
    fi

ARG GOMODCACHE
ARG GOCACHE
ARG GOOS
ARG GOARCH

COPY go.mod go.sum ./

RUN --mount=type=cache,target=${GOMODCACHE} go mod download

COPY . .

RUN --mount=type=cache,target=${GOCACHE} CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o zeus apps/olympus/zeus/main.go

FROM ubuntu:22.04 as doctl
RUN apt-get update && apt-get install -y wget
WORKDIR /app

RUN wget https://github.com/digitalocean/doctl/releases/download/v1.99.0/doctl-1.99.0-linux-amd64.tar.gz -O doctl.tar.gz
RUN tar xf doctl.tar.gz

FROM ubuntu:22.04 as gcp
RUN apt-get update && apt-get install -y wget
WORKDIR /app

RUN wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-429.0.0-linux-x86_64.tar.gz -O google-cloud-cli.tar.gz

FROM ubuntu:22.04 as aws_eks
RUN apt-get update && apt-get install -y wget curl
WORKDIR /app

ARG PLATFORM=Linux_amd64
RUN curl -sLO "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$PLATFORM.tar.gz"

# (Optional) Verify checksum
RUN curl -sL "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_checksums.txt" | grep $PLATFORM | sha256sum --check

RUN tar -xzf eksctl_$PLATFORM.tar.gz -C /tmp && rm eksctl_$PLATFORM.tar.gz

RUN mv /tmp/eksctl /app/eksctl

FROM ubuntu:22.04 as aws_cli
RUN apt-get update && apt-get install -y wget curl
WORKDIR /app

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"

FROM ubuntu:22.04 as cookbooks
RUN apt-get update && apt-get install -y git
WORKDIR /app

RUN git clone https://github.com/zeus-fyi/zeus.git --branch main --single-branch
RUN cp -r zeus/cookbooks .

# Stage 1: Download and unpack Helm
FROM ubuntu:22.04 as helm-installer
ARG HELM_VERSION=v3.0.0

# Install dependencies needed to fetch and unpack Helm
RUN apt-get update && apt-get install -y curl tar && rm -rf /var/lib/apt/lists/*

# Download and unpack Helm
WORKDIR /tmp
RUN curl -LO https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz \
    && tar -zxvf helm-${HELM_VERSION}-linux-amd64.tar.gz


FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates unzip less

COPY --from=builder /app/ca-certificate.crt /etc/ssl/certs
COPY --from=builder /app/ca.pem /etc/ssl/certs
COPY --from=builder /app/ca.key /etc/ssl/certs

COPY --from=aws_eks /app/eksctl /usr/local/bin
COPY --from=builder /app/zeus /usr/bin/zeus
COPY --from=doctl /app/doctl /usr/local/bin
COPY --from=gcp /app/google-cloud-cli.tar.gz /

COPY --from=aws_cli /app/awscliv2.zip /

# Copy the Helm binary from the helm-installer stage
COPY --from=helm-installer /tmp/linux-amd64/helm /usr/local/bin/helm

RUN unzip awscliv2.zip
RUN rm awscliv2.zip
RUN chmod +x ./aws/install

RUN ./aws/install

RUN tar xf google-cloud-cli.tar.gz
RUN rm google-cloud-cli.tar.gz
RUN ./google-cloud-sdk/install.sh
RUN /google-cloud-sdk/bin/gcloud components install gke-gcloud-auth-plugin
RUN chmod +x /google-cloud-sdk/bin/gcloud
RUN chmod +x /google-cloud-sdk/bin/gke-gcloud-auth-plugin

ENV PATH $PATH:/google-cloud-sdk/bin

COPY --from=cookbooks /app/cookbooks /etc/cookbooks
EXPOSE 9001