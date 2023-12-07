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

var DisplayConfigPath string
var DefaultConfig Config

type CoreSettings struct {
	BucketSize int `yaml:"bucket_size" env-default:"800"`
}

type BufferSize struct {
	Read  int `yaml:"read" env-default:"1024"`
	Write int `yaml:"write" env-default:"1024"`
}

type WebSocketConnector struct {
	Enable     bool       `yaml:"enable"`
	Address    string     `yaml:"address" env-default:"0.0.0.0"`
	Port       string     `yaml:"port" env-default:"8080"`
	BufferSize BufferSize `yaml:"buffer_size"`
}

type RestConnector struct {
	Enable  bool   `yaml:"enable" env-default:"True"`
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"31337"`
}

type GrpcConnector struct {
	Enable  bool   `yaml:"enable" env-default:"True"`
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"3137"`
}

type WebServer struct {
	Enable  bool   `yaml:"enable" env-default:"True"`
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"80"`
}

type Config struct {
	Env             string        `yaml:"env" env-default:"prod"`
	LogPath         string        `yaml:"log_path" env-default:"./logs/"`
	ShutdownTimeOut time.Duration `yaml:"shutdown_timeout" env-default:"5s"`

	CoreSettings       `yaml:"core_settings"`
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

	DisplayConfigPath = configPath
	DefaultConfig = cfg

	return &cfg
}
