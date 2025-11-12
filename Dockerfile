FROM golang:1.25.5

ARG ZEWI_VERSION=0.1.0

RUN mkdir -p /app/configs
RUN mkdir -p /app/var/logs
RUN apt-get update

WORKDIR /app

RUN curl -sL https://github.com/Clivern/Zewi/releases/download/v${ZEWI_VERSION}/zewi_Linux_x86_64.tar.gz | tar xz
RUN rm LICENSE
RUN rm README.md

COPY ./config.dist.yml /app/configs/

EXPOSE 8080

VOLUME /app/configs
VOLUME /app/var

RUN ./zewi version

CMD ["./zewi", "server", "-c", "/app/configs/config.dist.yml"]
