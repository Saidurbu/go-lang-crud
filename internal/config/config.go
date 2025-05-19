package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HTTPServer  `yaml:"http_server" env:"HTTP_SERVER" env-required:"true"`
	DBHost      string `yaml:"db_host" env:"DB_HOST" env-required:"true"`
	DBPort      string `yaml:"db_port" env:"DB_PORT" env-required:"true"`
	DBUser      string `yaml:"db_user" env:"DB_USER" env-required:"true"`
	DBPassword  string `yaml:"db_password" env:"DB_PASSWORD" env-required:"true"`
	DBName      string `yaml:"db_name" env:"DB_NAME" env-required:"true"`
}

func MustLoad() *Config {

	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	fmt.Println("Hello, welcome to crud api!", configPath)
	if configPath == "" {
		flags := flag.String("config", "", "Path to the config file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Please provide a config file path using the -config flag or set the CONFIG_PATH environment variable.")

		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}

	return &cfg
}
