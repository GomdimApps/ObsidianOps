#!/usr/bin/env bash

groupadd -f obeops
# Permissões das pasta de Backups
chgrp -R obeops /var/mine-backups/backup-server-bedrock/ /var/mine-backups/backup-server-java/

chgrp obeops /var/log/bedrock-backup.log /var/log/bedrock-update.log /var/log/java-server.log /var/log/bedrock-console.log
chgrp obeops /etc/obeops/bedrock-server.conf /etc/obeops/java-edition-server.conf /etc/obeops/mtools.conf
chmod 771 /var/log/bedrock-backup.log /var/log/bedrock-update.log /etc/obeops/bedrock-server.conf /etc/obeops/java-edition-server.conf /etc/obeops/mtools.conf
chmod 711 /usr/bin/mtools /usr/bin/bed-tools /usr/bin/console-bedrock /usr/bin/mtools-api

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

add_paths_to_bashrc() {
    local user_home=$1
    local bashrc_path="$user_home/.bashrc"
    local path_entry_bedrock="export PATH=\$PATH:/usr/bin/bedrock-tools/tools"
    local path_entry_java="export PATH=\$PATH:/usr/bin/java-tools"

    if ! grep -q "$path_entry_bedrock" "$bashrc_path"; then
        echo -e "\n# PATH do obeops" >> "$bashrc_path"
        echo "$path_entry_bedrock" >> "$bashrc_path"
    fi

    if ! grep -q "$path_entry_java" "$bashrc_path"; then
        echo "$path_entry_java" >> "$bashrc_path"
    fi
}

for user_home in /home/*; do
    if [ -d "$user_home" ]; then
        add_paths_to_bashrc "$user_home"
    fi
done

add_paths_to_bashrc "/root"
check_and_install
echo "Configuração finalizada."
clear
echo "Instalação do obeops finalizada"

exit 0