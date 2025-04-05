FROM golang:1.25.5

ARG ZEWI_VERSION=0.6.0

RUN mkdir -p /app/configs
RUN mkdir -p /app/var/logs
RUN apt-get update

WORKDIR /app

RUN curl -sL https://github.com/Clivern/Zewi/releases/download/v${ZEWI_VERSION}/zewi_Linux_x86_64.tar.gz | tar xz
RUN rm LICENSE
RUN rm README.md

COPY ./config.dist.yml /app/configs/
COPY ./entrypoint.sh /app/entrypoint.sh

EXPOSE 8080

VOLUME /app/configs
VOLUME /app/var

RUN chmod +x /app/entrypoint.sh
RUN chmod +x /app/zewi

CMD ["/app/entrypoint.sh", "api", "/app/configs/config.dist.yml"]
