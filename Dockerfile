# Etapa de construção
FROM golang:1.22-bullseye AS builder

WORKDIR /app

COPY . .

RUN apt-get update && \
    apt-get install -y make alien git curl && cd /app/ && make build-package

FROM debian:11-slim

EXPOSE 2308

COPY --from=builder /tmp/Build/APPS/obeops.deb /app/obeops.deb

RUN apt-get update && \
    apt-get install -y sudo && \
    sudo apt-get install -y dialog && \
    sudo apt-get install -y tar wget rsync tmux unzip lsof dialog && \
    sudo dpkg -i /app/obeops.deb

CMD ["/usr/bin/obeops", "--services", "-api-start"]