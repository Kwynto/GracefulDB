package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	ENV_DEV  = "dev"
	ENV_PROD = "prod"

	CONFIG_DEFAULT = "./config/default.yaml"
)

var DisplayConfigPath string
var DefaultConfig Config

type CoreSettings struct {
	Storage    string `yaml:"storage" env-default:"./data/"`
	BucketSize int64  `yaml:"bucket_size" env-default:"800"`
	// FreezeMode bool   `yaml:"freeze"`
	FriendlyMode bool `yaml:"friendly"`
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
	Enable  bool   `yaml:"enable"`
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"31337"`
}

type GrpcConnector struct {
	Enable  bool   `yaml:"enable"`
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"3137"`
}

type WebServer struct {
	Enable  bool   `yaml:"enable"`
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

func defaultConfig() Config {
	return Config{
		Env:             "test",
		LogPath:         "./logs/",
		ShutdownTimeOut: 5 * time.Second,
		CoreSettings: CoreSettings{
			Storage:      "./data/",
			BucketSize:   800,
			FriendlyMode: true,
		},
		WebSocketConnector: WebSocketConnector{
			Enable:  true,
			Address: "0.0.0.0",
			Port:    "8080",
			BufferSize: BufferSize{
				Read:  1024,
				Write: 1024,
			},
		},
		RestConnector: RestConnector{
			Enable:  true,
			Address: "0.0.0.0",
			Port:    "31337",
		},
		GrpcConnector: GrpcConnector{
			Enable:  true,
			Address: "0.0.0.0",
			Port:    "3137",
		},
		WebServer: WebServer{
			Enable:  true,
			Address: "0.0.0.0",
			Port:    "80",
		},
	}
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		configPath = CONFIG_DEFAULT
	}

	var cfg Config

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg = defaultConfig()
	} else if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		cfg = defaultConfig()
	}

	DisplayConfigPath = configPath
	DefaultConfig = cfg

	return &cfg
}
