package backup

import (
	"io/ioutil"
	"strings"
)

const ConfigFilePath = "/etc/obeops/bedrock-server.conf"

func GetServerDir() string {
	content, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		LogError("Erro ao ler o arquivo de configuração.")
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
