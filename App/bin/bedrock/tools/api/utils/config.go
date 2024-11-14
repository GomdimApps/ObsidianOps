package utils

import (
	"log"

	"github.com/GomdimApps/lcme"
)

const configPath = "/etc/mineservertools/mtools.conf"

type Config struct {
	ApiPort  int
	TokenApi string
}

func LoadConfig() (*Config, error) {
	config := Config{}
	err := lcme.ConfigRead(configPath, &config)
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
		return nil, err
	}
	return &config, nil
}
