# Etapa de construção
FROM golang:1.22-bullseye AS builder

WORKDIR /app

COPY . .

RUN bash /app/scripts/build/debian-11.sh && \
    sudo apt-get update && \
    sudo apt-get install -y make alien git curl && \
    make build-package

FROM debian:11-slim

COPY --from=builder /tmp/Build/APPS/obeops.deb /app/obeops.deb

RUN bash /app/scripts/build/debian-11.sh && \
    sudo apt-get update && \
    sudo apt-get install -y dialog && \
    sudo apt-get install -y tar wget rsync tmux unzip lsof dialog && \
    sudo dpkg -i /app/obeops.deb

CMD ["make", "build-package"]