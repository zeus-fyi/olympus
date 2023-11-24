#syntax=docker/dockerfile:1.4
FROM ubuntu:22.04 as src_builder

WORKDIR /app

COPY ./configs/secrets.tar.gz.age .
COPY ./configs/kube.tar.gz.age .
COPY ./configs/decrypt-and-setup.sh .

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates

RUN mkdir -p /etc/age
COPY --from=src_builder /app/secrets.tar.gz.age /etc/age
COPY --from=src_builder /app/kube.tar.gz.age /etc/age
COPY --from=src_builder /app/decrypt-and-setup.sh /etc

# Entry point script to decrypt and place files
RUN chmod +x /etc/decrypt-and-setup.sh

ENTRYPOINT ["/etc/decrypt-and-setup.sh"]
