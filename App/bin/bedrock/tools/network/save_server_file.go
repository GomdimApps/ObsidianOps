package network

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const ConfigFilePath = "/etc/obeops/bedrock-server.conf"

func SaveFile(filename string, content []byte, destDir string) error {
	updateDir := filepath.Join(destDir, "bedrock-update")
	if err := os.MkdirAll(updateDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(destDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return err
	}

	fmt.Printf("Arquivo salvo em %s, tamanho: %d bytes\n", filePath, len(content))

	if err := UnzipFile(filePath, updateDir); err != nil {
		return err
	}

	filesToRemove := []string{"server.properties", "permissions.json", "allowlist.json"}
	for _, file := range filesToRemove {
		err := os.Remove(filepath.Join(updateDir, file))
		if err != nil {
			return err
		}
	}

	cmd := "top -b -n 1 | grep -q 'bedrock' ; pkill bedrock_server || echo 'Server parado'"
	if err := executeCommand(cmd); err != nil {
		return fmt.Errorf("erro ao executar o comando para matar o processo: %v", err)
	}

	serverDir := GetServerDir()
	if serverDir == "" {
		return fmt.Errorf("diretório do servidor não encontrado")
	}

	if err := replaceServerFiles(updateDir, serverDir); err != nil {
		return fmt.Errorf("erro ao substituir arquivos do servidor: %v", err)
	}

	return nil
}

func executeCommand(cmd string) error {

	command := exec.Command("sh", "-c", cmd)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao executar o comando: %v, output: %s", err, output)
	}
	return nil
}

func GetServerDir() string {
	content, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo de configuração.")
		return ""
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "server-dir=") {
			return strings.Trim(strings.Split(line, "=")[1], "\"")
		}
	}
	return ""
}

func replaceServerFiles(updateDir, serverDir string) error {
	cmd := fmt.Sprintf("rsync -av %s/ %s/", updateDir, serverDir)
	if err := executeCommand(cmd); err != nil {
		return fmt.Errorf("erro ao executar o comando rsync: %v", err)
	}
	return nil
}
