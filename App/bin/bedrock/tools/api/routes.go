package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/GomdimApps/MineServerTools/App/bin/bedrock/tools/api/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	rateLimiters  = make(map[string]*rate.Limiter)
	ipBlacklist   = make(map[string]time.Time)
	ipLimit       = 25
	blockDuration = 40 * time.Second
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config, err := utils.LoadConfig()
		token_custom := config.TokenApi
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		token := c.GetHeader("Authorization")
		expectedToken := fmt.Sprintf("Bearer %s", token_custom)
		if token != expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if _, exists := rateLimiters[ip]; !exists {
			rateLimiters[ip] = rate.NewLimiter(3, 1)
		}

		if !rateLimiters[ip].Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ipBlockMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if blockUntil, blocked := ipBlacklist[ip]; blocked {
			if time.Now().Before(blockUntil) {
				c.JSON(http.StatusForbidden, gin.H{"error": "IP blocked due to excessive requests"})
				c.Abort()
				return
			}
			delete(ipBlacklist, ip)
		}

		if _, exists := rateLimiters[ip]; exists {
			if rateLimiters[ip].Allow() {
				countRequests(ip)
			}
		}

		c.Next()
	}
}

func countRequests(ip string) {
	if _, exists := ipBlacklist[ip]; !exists {
		ipBlacklist[ip] = time.Now().Add(blockDuration)
	}

	if len(ipBlacklist) > ipLimit {
		ipBlacklist[ip] = time.Now().Add(blockDuration)
	}
}

func SetupRoutes(r *gin.Engine) {
	authorized := r.Group("/", authMiddleware(), rateLimitMiddleware(), ipBlockMiddleware())
	{
		authorized.GET("/api/server/backup-bedrock/view/", viewBackup)
		authorized.GET("/api/server/status/view/", viewServerStatus)
		authorized.GET("/api/server/console-bedrock/log/", viewLog)
		authorized.POST("/api/server/backup-bedrock/task", performBackup)
		authorized.POST("/api/server/console-bedrock/start/", startServer)
	}
}

func viewBackup(c *gin.Context) {
	cmd := exec.Command("bed-tools", "--backup", "-j")
	stdout, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output := strings.TrimSpace(string(stdout))

	var backups []utils.Backup
	err = json.Unmarshal([]byte(output), &backups)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao parsear o JSON"})
		return
	}

	c.JSON(http.StatusOK, backups)
}

func viewServerStatus(c *gin.Context) {
	cmd := exec.Command("bed-tools", "--system", "-j")
	stdout, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output := strings.TrimSpace(string(stdout))

	var serverInfo utils.ServerInfo
	err = json.Unmarshal([]byte(output), &serverInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao parsear o JSON"})
		return
	}

	c.JSON(http.StatusOK, serverInfo)
}

func viewLog(c *gin.Context) {
	cmd := exec.Command("sh", "-c", `
		log_file="/var/log/bedrock-console.log"; 
		if [ -s "$log_file" ]; then 
			tail "$log_file" | sed -e 's/^[[:space:]]*//' -e 's/"/\\\"/g' | awk '{print "{\"file\": true, \"message\":\"" $0 "\"}"}'; 
		else 
			echo "{\"file\": false}"; 
		fi
	`)

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao executar o comando: %s", err.Error())})
		return
	}

	outputStr := strings.TrimSpace(string(output))

	if outputStr == "{\"file\": false}" {
		c.JSON(http.StatusOK, gin.H{"file": false})
		return
	}

	lines := strings.Split(outputStr, "\n")
	var logEntries []map[string]interface{}

	for _, line := range lines {
		if line != "" {
			var entry map[string]interface{}
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao parsear o JSON do log"})
				return
			}
			logEntries = append(logEntries, entry)
		}
	}

	c.JSON(http.StatusOK, logEntries)
}

func performBackup(c *gin.Context) {
	cmd := exec.Command("sh", "-c", `
		output=$(bed-tools --backup -b); 
		if echo "$output" | grep -q "Backup concluído"; then 
			echo '{"status": true, "message": "Seu backup foi realizado"}'; 
		else 
			echo '{"status": false, "message": "Erro backup"}'; 
		fi
	`)

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao executar o comando: %s", err.Error())})
		return
	}

	outputStr := strings.TrimSpace(string(output))

	c.JSON(http.StatusOK, outputStr)
}

func startServer(c *gin.Context) {
	cmd := exec.Command("sh", "-c", `
		output=$(console-bedrock --start); 
		if echo "$output" | grep -q "Iniciando o server"; then 
			echo '{"status": 200, "message": "Server iniciado"}'; 
		elif echo "$output" | grep -q "bedrock já foi iniciado"; then 
			echo '{"status": 467, "message": "O server já está ligado"}'; 
		else 
			echo '{"status": 500, "message": "Erro ao iniciar o server"}'; 
		fi
	`)

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao executar o comando: %s", err.Error())})
		return
	}

	outputStr := strings.TrimSpace(string(output))

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(outputStr), &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao parsear o JSON"})
		return
	}

	status := int(response["status"].(float64))
	c.JSON(status, response)
}
