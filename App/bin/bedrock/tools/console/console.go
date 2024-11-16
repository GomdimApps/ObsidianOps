package console

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/GomdimApps/ObsidianOps/App/bin/bedrock/tools/api/utils"
)

var sessionName = "server-bedrock"

func StartSession(dir string) error {
	if dir == "" {
		config, err := utils.LoadConfig()
		if err != nil {
			return err
		}
		dir = config.ServerDir
	} else {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("diretório especificado '%s' não existe", dir)
		}
	}
	if _, err := exec.Command("sh", "-c", "tmux has-session -t "+sessionName).Output(); err == nil {
		fmt.Println("O server bedrock já foi iniciado!")
		return nil
	}
	cmd := fmt.Sprintf("cd %s ; truncate -s 0 /var/log/bedrock-console.log ; LD_LIBRARY_PATH=. ./bedrock_server | tee -a /var/log/bedrock-console.log", dir)
	if _, err := exec.Command("sh", "-c", "tmux new-session -d -s "+sessionName+" "+cmd).Output(); err != nil {
		return fmt.Errorf("erro ao iniciar a sessão tmux: %v", err)
	}
	fmt.Printf("Iniciando o server bedrock no diretório '%s'...\n", dir)
	return nil
}

func StopSession() error {
	if _, err := exec.Command("sh", "-c", "tmux has-session -t "+sessionName).Output(); err == nil {
		if _, err := exec.Command("sh", "-c", "tmux kill-session -t "+sessionName).Output(); err != nil {
			return fmt.Errorf("erro ao parar a sessão tmux: %v", err)
		}
		if _, err := exec.Command("sh", "-c", "truncate -s 0 /var/log/bedrock-console.log").Output(); err != nil {
			return fmt.Errorf("erro ao truncar o arquivo de log: %v", err)
		}
		fmt.Println("Parando Server Bedrock")
	} else {
		fmt.Println("O server bedrock não está ligado.")
	}
	return nil
}

func ConnectSession() error {
	if _, err := exec.Command("sh", "-c", "tmux has-session -t "+sessionName).Output(); err == nil {
		if _, err := exec.Command("sh", "-c", "tmux attach -t "+sessionName).Output(); err != nil {
			return fmt.Errorf("erro ao conectar à sessão tmux: %v", err)
		}
	} else {
		fmt.Println("O server bedrock não está ligado.")
	}
	return nil
}

func SendCommand(command string) error {
	if _, err := exec.Command("sh", "-c", "tmux has-session -t "+sessionName).Output(); err == nil {
		if _, err := exec.Command("sh", "-c", "tmux send-keys -t "+sessionName+" "+command+" C-m").Output(); err != nil {
			return fmt.Errorf("erro ao enviar comando: %v", err)
		}
		fmt.Printf("Comando '%s' enviado para o server bedrock.\n", command)
	} else {
		fmt.Println("O server bedrock não está ligado.")
	}
	return nil
}
