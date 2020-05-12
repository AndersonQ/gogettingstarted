package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
)

type Config struct {
	AppName string `env:"APP_NAME" envDefault:"go-gettingstarted"`
	Env     string `env:"ENV" envDefault:"env not set"`

	LogLevel  string `env:"LOG_LEVEL" envDefault:"debug"`
	LogOutput string `env:"LOG_OUTPUT" envDefault:"console"`
}

func ParseManual() (Config, error) {
	key := "APP_NAME"
	val, ok := os.LookupEnv(key)
	if !ok {
		return Config{}, fmt.Errorf("environment variable %s not defined", key)
	}

	return Config{AppName: val}, nil
}

func Parse() (Config, error) {
	cfg := Config{}
	err := env.Parse(&Config{})
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return cfg, nil
}

// Logger returns a initialised zerolog.Logger
func (c Config) Logger() zerolog.Logger {
	logLevelOk := true
	logLevel, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
		logLevelOk = false
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimestampFieldName = "timestamp"

	host, _ := os.Hostname()
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("application", c.AppName).
		Str("host", host).
		Str("environment", c.Env).
		Logger()

	if strings.ToUpper(c.LogOutput) == "CONSOLE" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if !logLevelOk {
		logger.Warn().Err(err).Msgf("%s is not a valid zerolog log level, defaulting to info", c.LogLevel)
	}

	return logger
}
