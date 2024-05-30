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

var SDisplayConfigPath string
var StDefaultConfig TConfig

type TCoreSettings struct {
	Storage      string `yaml:"storage" env-default:"./data/"`
	BucketSize   int64  `yaml:"bucket_size" env-default:"800"`
	FriendlyMode bool   `yaml:"friendly"`
}

type TBufferSize struct {
	Read  int `yaml:"read" env-default:"1024"`
	Write int `yaml:"write" env-default:"1024"`
}

type TWebSocketConnector struct {
	Enable     bool        `yaml:"enable"`
	Address    string      `yaml:"address" env-default:"0.0.0.0"`
	Port       string      `yaml:"port" env-default:"8080"`
	BufferSize TBufferSize `yaml:"buffer_size"`
}

type TRestConnector struct {
	Enable  bool   `yaml:"enable"`
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"31337"`
}

type TGrpcConnector struct {
	Enable  bool   `yaml:"enable"`
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"3137"`
}

type TWebServer struct {
	Enable  bool   `yaml:"enable"`
	Address string `yaml:"address" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"80"`
}

type TConfig struct {
	Env             string        `yaml:"env" env-default:"prod"`
	LogPath         string        `yaml:"log_path" env-default:"./logs/"`
	ShutdownTimeOut time.Duration `yaml:"shutdown_timeout" env-default:"5s"`

	CoreSettings       TCoreSettings       `yaml:"core_settings"`
	WebSocketConnector TWebSocketConnector `yaml:"websocket_connector"`
	RestConnector      TRestConnector      `yaml:"rest_connector"`
	GrpcConnector      TGrpcConnector      `yaml:"grpc_connector"`
	WebServer          TWebServer          `yaml:"web_server"`
}

func defaultConfig() TConfig {
	return TConfig{
		Env:             "test",
		LogPath:         "./logs",
		ShutdownTimeOut: 5 * time.Second,
		CoreSettings: TCoreSettings{
			Storage:      "./data",
			BucketSize:   800,
			FriendlyMode: true,
		},
		WebSocketConnector: TWebSocketConnector{
			Enable:  true,
			Address: "0.0.0.0",
			Port:    "8080",
			BufferSize: TBufferSize{
				Read:  1024,
				Write: 1024,
			},
		},
		RestConnector: TRestConnector{
			Enable:  true,
			Address: "0.0.0.0",
			Port:    "31337",
		},
		GrpcConnector: TGrpcConnector{
			Enable:  true,
			Address: "0.0.0.0",
			Port:    "3137",
		},
		WebServer: TWebServer{
			Enable:  true,
			Address: "0.0.0.0",
			Port:    "80",
		},
	}
}

func MustLoad(sConfigPath string) *TConfig {
	if sConfigPath == "" {
		sConfigPath = CONFIG_DEFAULT
	}

	var stCfg TConfig

	// check if file exists
	if _, err := os.Stat(sConfigPath); os.IsNotExist(err) {
		stCfg = defaultConfig()
	} else if err := cleanenv.ReadConfig(sConfigPath, &stCfg); err != nil {
		stCfg = defaultConfig()
	}

	SDisplayConfigPath = sConfigPath
	StDefaultConfig = stCfg

	return &stCfg
}
