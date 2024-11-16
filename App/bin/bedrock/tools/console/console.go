package console

import (
	"fmt"
	"os"

	"github.com/GomdimApps/ObsidianOps/App/bin/bedrock/tools/api/utils"

	"github.com/GomdimApps/lcme"
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
	if _, err := lcme.Shell("tmux has-session -t " + sessionName); err == nil {
		fmt.Println("O server bedrock já foi iniciado!")
		return nil
	}
	cmd := fmt.Sprintf("cd %s ; truncate -s 0 /var/log/bedrock-console.log ; LD_LIBRARY_PATH=. ./bedrock_server | tee -a /var/log/bedrock-console.log", dir)
	if _, err := lcme.Shell("tmux new-session -d -s " + sessionName + " " + cmd); err != nil {
		return fmt.Errorf("erro ao iniciar a sessão tmux: %v", err)
	}
	fmt.Printf("Iniciando o server bedrock no diretório '%s'...\n", dir)
	return nil
}

func StopSession() error {
	if _, err := lcme.Shell("tmux has-session -t " + sessionName); err == nil {
		if _, err := lcme.Shell("tmux kill-session -t " + sessionName); err != nil {
			return fmt.Errorf("erro ao parar a sessão tmux: %v", err)
		}
		if _, err := lcme.Shell("truncate -s 0 /var/log/bedrock-console.log"); err != nil {
			return fmt.Errorf("erro ao truncar o arquivo de log: %v", err)
		}
		fmt.Println("Parando Server Bedrock")
	} else {
		fmt.Println("O server bedrock não está ligado.")
	}
	return nil
}

func ConnectSession() error {
	if _, err := lcme.Shell("tmux has-session -t " + sessionName); err == nil {
		if _, err := lcme.Shell("tmux attach -t " + sessionName); err != nil {
			return fmt.Errorf("erro ao conectar à sessão tmux: %v", err)
		}
	} else {
		fmt.Println("O server bedrock não está ligado.")
	}
	return nil
}

func SendCommand(command string) error {
	if _, err := lcme.Shell("tmux has-session -t " + sessionName); err == nil {
		if _, err := lcme.Shell("tmux send-keys -t " + sessionName + " " + command + " C-m"); err != nil {
			return fmt.Errorf("erro ao enviar comando: %v", err)
		}
		fmt.Printf("Comando '%s' enviado para o server bedrock.\n", command)
	} else {
		fmt.Println("O server bedrock não está ligado.")
	}
	return nil
}
