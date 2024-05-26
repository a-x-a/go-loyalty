package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

type (
	ServiceConfig struct {
		// ListenAddress - адрес и порт запуска сервиса.
		RunAddress string `env:"RUN_ADDRESS"`
		// DatabaseURI - адрес подключения к базе данных.
		DatabaseURI string `env:"DATABASE_URI"`
		// AccrualSystemAddress - адрес системы расчёта начислений.
		AccrualSystemAddress string `env:"ACCRUEAL_SYSTEM_ADDRESS"`
		// Secret - секретная фраза.
		Secret string
	}
)

func NewServiceConfig() ServiceConfig {
	cfg := ServiceConfig{
		RunAddress:           "localhost:8080",
		DatabaseURI:          "",
		AccrualSystemAddress: "",
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Использование:\n")
		flag.PrintDefaults()
	}

	if flag.Lookup("a") == nil {
		flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "адрес и порт запуска сервиса")
	}

	if flag.Lookup("d") == nil {
		flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "адрес подключения к базе данных")
	}

	if flag.Lookup("r") == nil {
		flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "адрес системы расчёта начислений")
	}

	flag.Parse()

	_ = env.Parse(&cfg)

	return cfg
}
