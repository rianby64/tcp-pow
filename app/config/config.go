package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address             string `env:"ADDRESS,required"`
	LeadingBits         uint   `env:"LEADING_BITS,required"`
	SaltSize            uint   `env:"SALT_SIZE,required"`
	KeySize             int    `env:"KEY_SIZE" envDefault:"64"`
	ProcessTimeoutSecs  int    `env:"PROCESS_TIMEOUT_SECS" envDefault:"2"`
	ShutdownTimeoutSecs int    `env:"SHUTDOWN_TIMEOUT_SECS" envDefault:"30"`
	QuotesPath          string `env:"QUOTES_PATH,required"`
}

func CreateLogAndConfig() (*log.Logger, *Config) {
	cfg := &Config{}
	log := log.New(os.Stderr, "server-tcp-pow", log.LstdFlags)

	if err := env.Parse(cfg); err != nil {
		log.Panicf("env.Parse(&cfg): %v", err)
	}

	log.Printf("ADDRESS=%v", cfg.Address)
	log.Printf("LEADING_BITS=%v", cfg.LeadingBits)
	log.Printf("SALT_SIZE=%v", cfg.SaltSize)
	log.Printf("KEY_SIZE=%v", cfg.KeySize)
	log.Printf("PROCESS_TIMEOUT_SECS=%v", cfg.ProcessTimeoutSecs)
	log.Printf("SHUTDOWN_TIMEOUT_SECS=%v", cfg.ShutdownTimeoutSecs)
	log.Printf("QUOTES_PATH=%v", cfg.QuotesPath)

	return log, cfg
}
