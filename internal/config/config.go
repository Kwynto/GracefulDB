package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"

	configDefault = "./config/default.yaml"
)

type CoreSettings struct {
	BucketSize int `yaml:"bucket_size" env-default:"800" env-required:"true"`
}

type SocketConnector struct {
	Enable bool `yaml:"enable" env-default:"true" env-required:"true"`
}

type WebSocketConnector struct {
	Enable bool `yaml:"enable" env-default:"true" env-required:"true"`
}

type RestConnector struct {
	Enable  bool   `yaml:"enable" env-default:"true" env-required:"true"`
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"31337"`
}

type GrpcConnector struct {
	Enable  bool   `yaml:"enable" env-default:"true" env-required:"true"`
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"8080"`
}

type WebServer struct {
	Enable  bool   `yaml:"enable" env-default:"true" env-required:"true"`
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"80"`
}

type Config struct {
	Env             string        `yaml:"env" env-default:"prod"`
	LogPath         string        `yaml:"log_path" env-default:"./logs/"`
	ShutdownTimeOut time.Duration `yaml:"shutdown_timeout" env-default:"5s"`

	CoreSettings       `yaml:"core_settings"`
	SocketConnector    `yaml:"socket_connector"`
	WebSocketConnector `yaml:"websocket_connector"`
	RestConnector      `yaml:"rest_connector"`
	GrpcConnector      `yaml:"grpc_connector"`
	WebServer          `yaml:"web_server"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		configPath = configDefault
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", configPath)
	}

	return &cfg
}
