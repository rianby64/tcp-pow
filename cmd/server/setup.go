package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rianby64/tcp-pow/server"
)

const (
	shutdownTimeout = time.Second * 30
)

type Shutdown interface {
	Shutdown(ctx context.Context) error
}

type Config struct {
	Address     string `env:"ADDRESS,required"`
	LeadingBits uint   `env:"LEADING_BITS,required"`
	SaltSize    uint   `env:"SALT_SIZE,required"`
	KeySize     int    `env:"KEY_SIZE" envDefault:"64"`
	QuotesPath  string `env:"QUOTES_PATH,required"`
}

func createLogAndConfig() (*log.Logger, *Config) {
	cfg := &Config{}
	log := log.New(os.Stderr, "server-tcp-pow", log.LstdFlags)

	if err := env.Parse(cfg); err != nil {
		log.Panicf("env.Parse(&cfg): %v", err)
	}

	log.Printf("ADDRESS=%v", cfg.Address)
	log.Printf("LEADING_BITS=%v", cfg.LeadingBits)
	log.Printf("SALT_SIZE=%v", cfg.SaltSize)
	log.Printf("QUOTES_PATH=%v", cfg.QuotesPath)

	return log, cfg
}

func waitSignalShutdown() {
	signalShutdown := make(chan os.Signal, 1)

	signal.Notify(signalShutdown, syscall.SIGINT, syscall.SIGTERM)

	<-signalShutdown
}

func shutdown(log *log.Logger, items ...Shutdown) {
	log.Printf("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	defer cancel()

	for _, item := range items {
		if err := item.Shutdown(ctx); err != nil {
			log.Printf("err := item.Shutdown(ctx): %v", err)
		}
	}

	log.Printf("shut down completed")
}

func listen(server *server.Server, log *log.Logger, address string) {
	log.Printf("Listening at %s", address)

	err := server.Listen(address)

	if errors.Is(err, net.ErrClosed) {
		log.Printf("Listener closed")

		return
	}

	if err != nil {
		shutdown(log, server)

		log.Panicf("server.Listen(%s): %v", address, err)
	}
}
