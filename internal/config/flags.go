package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v10"
)

func parseFlags(config config) config {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Флаг -a=<ЗНАЧЕНИЕ> адрес и порт запуска сервиса (по умолчанию localhost:8080).
	runAddr := flag.String("a", "localhost:8080", "service address and port")

	// Флаг -r=<ЗНАЧЕНИЕ> адрес системы расчёта начислений (по умолчанию localhost:8081).
	accrualAddr := flag.String("r", "localhost:8081", "address:port of the accrual calculation system")

	// Флаг -d=<ЗНАЧЕНИЕ> адрес подключения к базе данных
	databaseURI := flag.String("d", "", "database URI")

	flag.Parse()

	return config.
		SetRunAddr(*runAddr).
		SetAccrualAddr(*accrualAddr).
		SetDatabaseURI(*databaseURI)
}

func parseEnvs(config config) config {
	var cfg struct {
		RunAddr     string `env:"RUN_ADDRESS"`
		AccrualAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
		DatabaseURI string `env:"DATABASE_URI"`
	}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	if _, exists := os.LookupEnv("RUN_ADDRESS"); exists {
		config = config.SetRunAddr(cfg.RunAddr)
	}

	if _, exists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); exists {
		config = config.SetAccrualAddr(cfg.AccrualAddr)
	}

	if _, exists := os.LookupEnv("DATABASE_URI"); exists {
		config = config.SetDatabaseURI(cfg.DatabaseURI)
	}

	return config
}
