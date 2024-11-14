#!/bin/bash

apt-get update && apt-get install -y sudo

sudo rm -rf /etc/apt/sources.list

sudo tee /etc/apt/sources.list <<EOF
deb http://ftp.br.debian.org/debian bullseye main
deb http://ftp.br.debian.org/debian-security bullseye-security main
deb http://ftp.br.debian.org/debian bullseye-updates main
EOF

sudo apt-get update && sudo apt-get install -y make alien git curl tar wget rsync tmux unzip lsof dialog