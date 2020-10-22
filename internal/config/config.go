package config

import (
	"github.com/namsral/flag"
)

const (
	DefaultHTTPPort = ":8080"
	DefaultGRPCPort = ":9090"
	EnvDev          = "development"
	EnvProd         = "production"
	LogLevelInfo    = "info"
	LogLevelDebug   = "debug"
)

type Config struct {
	// App Config
	AppName  string // brackets by default
	Env      string // development or production
	LogLevel string // debug / info

	// Server Config
	HTTPPort string
	GRPCPort string

	// Usage
	Usage bool
}

func MakeConfig() (Config, error) {
	cfg := Config{}
	// App Config
	flag.StringVar(&cfg.AppName, "APP_NAME", "brackets", "Application name")
	flag.StringVar(&cfg.Env, "ENVIRONMENT", EnvProd, "Application env")
	flag.StringVar(&cfg.LogLevel, "LOG_LEVEL", LogLevelInfo, "Log level")
	// Server Config
	flag.StringVar(&cfg.HTTPPort, "HTTP_PORT", DefaultHTTPPort, "Http port")
	flag.StringVar(&cfg.GRPCPort, "GRPC_PORT", DefaultGRPCPort, "GRPC port")

	// Usage
	flag.BoolVar(&cfg.Usage, "USAGE", false, "Show usage")

	flag.Parse()

	return cfg, nil
}

func Usage() {
	flag.Usage()
}
