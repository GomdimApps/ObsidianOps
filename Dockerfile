# Etapa de construção
FROM golang:1.22-bullseye AS builder

WORKDIR /app

COPY . .

RUN bash /app/scripts/build/debian-11.sh

CMD ["make", "build-package"]