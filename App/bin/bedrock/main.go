package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GomdimApps/ObsidianOps/App/bin/bedrock/tools/api"
	"github.com/GomdimApps/ObsidianOps/App/bin/bedrock/tools/api/utils"
	"github.com/GomdimApps/ObsidianOps/App/bin/bedrock/tools/backup"
	"github.com/GomdimApps/ObsidianOps/App/bin/bedrock/tools/network"
	"github.com/GomdimApps/ObsidianOps/App/bin/bedrock/tools/system"
	"github.com/gin-gonic/gin"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Erro: Argumentos insuficientes.")
		fmt.Println("Uso: obeops [--backup|--schedule|--view]")
		fmt.Println("Use --help para mais informações.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--backup":
		switch os.Args[2] {
		case "-b":
			backup.Backup()
		case "-j":
			jsonOutput, err := backup.ViewBackupJson()
			if err != nil {
				log.Fatalf("Erro: %v", err)
			}
			fmt.Println(jsonOutput)
		case "-s":
			backup.ScheduleTasks()
		case "-v":
			backup.ViewBackup()
		default:
			fmt.Println("Erro: Opção inválida para --backup.")
			fmt.Println("Uso: obeops --backup [-b|-j|-s|-v]")
			fmt.Println("Use --help para mais informações.")
			os.Exit(1)
		}
	case "--system":
		switch os.Args[2] {
		case "-b":
			system.PrintStatus()
		case "-j":
			jsonOutputSystem, err := system.GetStatusJSON()
			if err != nil {
				log.Fatalf("Erro: %v", err)
			}
			fmt.Println(jsonOutputSystem)
		default:
			fmt.Println("Erro: Opção inválida para --schedule.")
			fmt.Println("Uso: obeops --schedule [-b]")
			fmt.Println("Use --help para mais informações.")
			os.Exit(1)
		}
	case "--server":
		switch os.Args[2] {
		case "-u":
			err := network.RunDownloadProcess()
			if err != nil {
				fmt.Printf("Erro: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Atualização do servidor completada")
		case "-n":
			if len(os.Args) < 4 {
				fmt.Println("Erro: O diretório para o novo servidor deve ser especificado.")
				os.Exit(1)
			}
			destDir := os.Args[3]
			err := network.SetupNewServer(destDir)
			if err != nil {
				fmt.Printf("Erro: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Novo servidor instalado")
		default:
			fmt.Println("Erro: Opção inválida para --view.")
			fmt.Println("Uso: obeops --view [-b|-j]")
			fmt.Println("Use --help para mais informações.")
			os.Exit(1)
		}

	case "--services":
		switch os.Args[2] {
		case "-api-start":
			r := gin.Default()
			config, err := utils.LoadConfig()
			apiPort := config.ApiPort
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			api.SetupRoutes(r)

			r.Run(fmt.Sprintf(":%d", apiPort))

		default:
			fmt.Println("Erro: Opção inválida para --services.")
			fmt.Println("Uso: obeops --services -api-start")
			fmt.Println("Use --help para mais informações.")
			os.Exit(1)
		}
	case "--help":
		fmt.Println("Ajuda: Uso do programa.")
		fmt.Println("obeops [--backup|--schedule|--view] [-b|-j]")
		fmt.Println("  --backup  : Realiza um backup.")
		fmt.Println("  --schedule: Agenda uma tarefa.")
		fmt.Println("  --view    : Visualiza um backup.")
		fmt.Println("  -b        : Opção para tipo B.")
		fmt.Println("  -j        : Opção para tipo J.")

	default:
		fmt.Println("Erro: Comando inválido.")
		fmt.Println("Uso: obeops [--backup|--schedule|--view] [-b|-j]")
		fmt.Println("Use --help para mais informações.")
		os.Exit(1)
	}

}
