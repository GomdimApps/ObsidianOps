package system

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

const ConfigFilePath = "/etc/obeops/bedrock-server.conf"

func GetServerDir() string {
	content, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo de configuração do Bedrock")
		return ""
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "server-dir=") {
			return strings.Trim(strings.Split(line, "=")[1], "\"")
		}
	}
	return ""
}

func DirSize(path string) (uint64, error) {
	var size uint64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += uint64(info.Size())
		}
		return nil
	})
	return size, err
}

func GetDiskUsage(path string) (totalBytes, freeBytes, usedBytes uint64, usedPercent float64) {
	usage, _ := disk.Usage(path)
	return usage.Total, usage.Free, usage.Used, usage.UsedPercent
}

func GetServerIP() (string, error) {
	cmd := exec.Command("sh", "-c", "ip route get 8.8.8.8 | awk '{print $7}'")
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}

func GetServerStatus() (string, error) {
	cmd := exec.Command("sh", "-c", "top -b -n 1 | grep -q \"bedrock\" && echo \"Operando\" || echo \"Inativo\"")
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}

func GetPortsUsed() (string, error) {
	cmd := exec.Command("sh", "-c", `output=$(lsof -i -P -n | grep bedrock | awk '/IPv4/ {print "IPv4: " $9} /IPv6/ {print "IPv6: " $9}'); if [ -z "$output" ]; then echo "Servidor inativo"; else echo "$output"; fi`)
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}

func GetStatusJSON() (string, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "", fmt.Errorf("erro ao obter a porcentagem de uso da CPU: %v", err)
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		return "", fmt.Errorf("erro ao obter o uso da memória: %v", err)
	}

	totalDisk, freeDisk, usedDisk, usedPercent := GetDiskUsage("/")
	serverDir := GetServerDir()
	if serverDir == "" {
		return "", fmt.Errorf("erro ao obter o diretório do servidor")
	}

	serverDirSize, err := DirSize(serverDir)
	if err != nil {
		return "", fmt.Errorf("erro ao obter o tamanho do diretório do servidor: %v", err)
	}

	serverStatus, err := GetServerStatus()
	if err != nil {
		return "", fmt.Errorf("erro ao verificar o status do servidor: %v", err)
	}

	serverIP, err := GetServerIP()
	if err != nil {
		return "", fmt.Errorf("erro ao obter o IP do servidor: %v", err)
	}

	portsUsed, err := GetPortsUsed()
	if err != nil {
		return "", fmt.Errorf("erro ao obter as portas utilizadas: %v", err)
	}

	status := map[string]interface{}{
		"CPUUsage":    cpuPercent[0],
		"MemoryUsage": memory.UsedPercent,
		"TotalMemory": memory.Total,
		"FreeMemory":  memory.Free,
		"TotalDisk":   totalDisk,
		"FreeDisk":    freeDisk,
		"UsedDisk": map[string]interface{}{
			"Size":       usedDisk,
			"Percentage": usedPercent,
		},
		"ServerSize":   serverDirSize,
		"ServerStatus": serverStatus,
		"ServerIP":     serverIP,
		"PortsUsed":    portsUsed,
	}

	jsonData, err := json.Marshal(status)
	if err != nil {
		return "", fmt.Errorf("erro ao criar o JSON de status: %v", err)
	}

	return string(jsonData), nil
}

func PrintStatus() {
	cpuPercent, _ := cpu.Percent(time.Second, false)
	memory, _ := mem.VirtualMemory()
	totalDisk, freeDisk, usedDisk, usedPercent := GetDiskUsage("/")
	serverDir := GetServerDir()

	if memory.UsedPercent >= 85 {
		fmt.Println("##### AVISO: a utilização da memória é superior a 85%")
	}

	serverDirSize, _ := DirSize(serverDir)
	serverStatus, _ := GetServerStatus()
	serverIP, _ := GetServerIP()
	portsUsed, _ := GetPortsUsed()

	fmt.Printf("\nUtilização da CPU: %.2f%%\n", cpuPercent[0])
	fmt.Printf("Utilização da memória: %.2f%%\n\n", memory.UsedPercent)
	fmt.Printf("Memória Total: %v MB\n", memory.Total/1024/1024)
	fmt.Printf("Memória livre: %v MB\n", memory.Free/1024/1024)
	fmt.Printf("Total de Disco: %v GB\n", totalDisk/1024/1024/1024)
	fmt.Printf("Disco Livre: %v GB\n", freeDisk/1024/1024/1024)
	fmt.Printf("Disco Usado: %v GB (%.2f%%)\n\n", usedDisk/1024/1024/1024, usedPercent)

	if serverDirSize < 1024*1024*1024 {
		fmt.Printf("Tamanho do servidor: %v MB\n", serverDirSize/1024/1024)
	} else {
		fmt.Printf("Tamanho do servidor: %.2f GB\n", float64(serverDirSize)/(1024*1024*1024))
	}

	fmt.Printf("Status do servidor: %s\n", serverStatus)
	fmt.Printf("IP do servidor: %s\n", serverIP)
	fmt.Printf("Portas utilizadas:\n%s\n\n", portsUsed)
}
