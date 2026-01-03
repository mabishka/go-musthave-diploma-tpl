package config

import (
	"flag"
	"os"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"go.uber.org/zap"
)

type Config struct {
	RunAddress           string
	DatabaseURL          string
	AccuralSystemAddress string
	LogLevel             string
}

const (
	envRunAddress           = "RUN_ADDRESS"
	envDatabaseURL          = "DATABASE_URL"
	envAccuralSystemAddress = "ACCURAL_SYSTEM_ADDRESS"
	envLogLevel             = "LOG_LEVEL"

	flagRunAddress           = "a"
	flagDatabaseURL          = "d"
	flagAccuralSystemAddress = "r"
	flagLogLevel             = "l"

	descRunAddress           = "адрес и порт запуска сервиса"
	descDatabaseURL          = "адрес подключения к базе данных"
	descAccuralSystemAddress = "адрес системы расчёта начислений"
	descLogLevel             = "уровень логирования"
)

func New() *Config {
	return &Config{
		RunAddress:           "localhost:8000",
		DatabaseURL:          "",
		AccuralSystemAddress: "localhost:8080",
		LogLevel:             "info",
	}
}

func (c *Config) Load() {
	a := setValue(envRunAddress, flagRunAddress, descRunAddress)
	d := setValue(envDatabaseURL, flagDatabaseURL, descDatabaseURL)
	r := setValue(envAccuralSystemAddress, flagAccuralSystemAddress, descAccuralSystemAddress)
	l := setValue(envLogLevel, flagLogLevel, descLogLevel)

	flag.Parse()

	c.RunAddress = *a
	c.DatabaseURL = *d
	c.AccuralSystemAddress = *r
	if *l != "" {
		c.LogLevel = *l
	}
}

func setValue(envName, flagName, desc string) *string {
	flagValue := flag.String(flagName, "", desc)
	if envValue, ok := os.LookupEnv(envName); ok && envValue != "" {
		logger.Log().Info("set env value", zap.String(envName, envValue))
		return &envValue
	}
	logger.Log().Info("set flag value", zap.String(envName, flagName))
	return flagValue
}
