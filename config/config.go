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
	Database   *DatabaseConfig   `yaml:"database"`
	Redis      *RedisConfig      `yaml:"redis"`
	GRPC       *GRPCConfig       `yaml:"grpc"`
	Client     *ClientConfig     `yaml:"client"`
}

type HTTPServerConfig struct {
	Address string `yaml:"address"`
}

type DatabaseConfig struct {
	ConnectionString  string `yaml:"connection_string"`
	DbType            string `yaml:"db_type"`
	MigrationFilePath string `yaml:"migration_file_path"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
	PoolSize int    `yaml:"pool-size"`
}

type GRPCConfig struct {
	Port string `yaml:"port"`
}

type ClientConfig struct {
	UserService string `yaml:"user_service"`
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
