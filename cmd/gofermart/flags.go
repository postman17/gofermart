package main

import (
	"flag"
	"log/slog"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	RunAddr              string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func parseFlags() Config {
	var (
		runAddr              string
		databaseURI          string
		accrualSystemAddress string
	)
	flag.StringVar(&runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(
		&databaseURI,
		"d",
		"postgres://gofermart_user:gofermart_user_password@localhost:5432/gofermart",
		"database address",
	)
	flag.StringVar(&accrualSystemAddress, "r", "localhost:8081", "address and port accrual server")
	flag.Parse()

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		slog.Error("Error parse envs:", "error", err)
	}

	if cfg.RunAddr != "" {
		runAddr = cfg.RunAddr
	}
	if cfg.DatabaseURI != "" {
		databaseURI = cfg.DatabaseURI
	}
	if cfg.AccrualSystemAddress != "" {
		accrualSystemAddress = cfg.AccrualSystemAddress
	}

	return Config{
		RunAddr:              runAddr,
		DatabaseURI:          databaseURI,
		AccrualSystemAddress: accrualSystemAddress,
	}
}
