#!/usr/bin/env bash

groupadd -f obeops
# Permissões das pasta de Backups
mkdir -p /var/mine-backups/backup-server-bedrock/
chgrp -R obeops /var/mine-backups/backup-server-bedrock/ 

chgrp obeops /var/log/bedrock-backup.log /var/log/bedrock-update.log /var/log/java-server.log /var/log/bedrock-console.log
chgrp obeops /etc/obeops/bedrock-server.conf /etc/obeops/java-edition-server.conf /etc/obeops/mtools.conf
chmod 771 /var/log/bedrock-backup.log /var/log/bedrock-update.log /etc/obeops/bedrock-server.conf /etc/obeops/java-edition-server.conf /etc/obeops/mtools.conf
chmod 711 /usr/bin/obeops /usr/bin/ob-tools /usr/bin/console-bedrock 

dependencies=(tar wget rsync tmux unzip lsof)

check_and_install() {
    echo "Verificando dependências..."
    sudo apt update
    for dep in "${dependencies[@]}"; do
        if ! dpkg -l | grep -q "$dep"; then
            echo "$dep não está instalado. Instalando..."
            sudo apt install "$dep" -y
        else
            echo "$dep já está instalado."
        fi
    done
    echo "Todas as dependências foram instaladas..."
}

check_and_install
echo "Configuração finalizada."
clear
echo "Instalação do obeops finalizada"

exit 0