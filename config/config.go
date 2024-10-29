package config

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerName string            `yaml:"service_name"`
	HTTPServer *HTTPServerConfig `yaml:"http_server"`
}

type HTTPServerConfig struct {
	Address string `yaml:"address"`
}

func LoadConfig() *Config {

	configPath := os.Getenv("CONFIG_PATH")

	// check configPath in env
	if configPath == "" {
		// get configPath from command
		flagConfigPath := flag.String("config", "", "Path to the configuration file")
		flag.Parse()

		if *flagConfigPath == "" {
			log.Fatal("Configuration file path is required. Use --config=<path>")
		}
		configPath = *flagConfigPath
	}

	var cfg Config

	yamlData, err := os.ReadFile(configPath)

	if err != nil {
		log.Fatal("Error while reading config file ", err)
	}

	yaml.Unmarshal(yamlData, &cfg)

	return &cfg
}
