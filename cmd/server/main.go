package main

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/catalinc/hashcash"
	"github.com/rianby64/tcp-pow/book"
	"github.com/rianby64/tcp-pow/handler/quotes"
	"github.com/rianby64/tcp-pow/handler/verify"
	"github.com/rianby64/tcp-pow/server"
)

type Config struct {
	Address      string `env:"ADDRESS,required"`
	LeadingBytes uint   `env:"LEADING_BYTES,required"`
	SaltSize     uint   `env:"SALT_SIZE,required"`
	QuotesPath   string `env:"QUOTES_PATH,required"`
}

func main() {
	cfg := &Config{}
	log := log.New(os.Stderr, "server-tcp-pow", log.LstdFlags)

	if err := env.Parse(cfg); err != nil {
		log.Panicf("env.Parse(&cfg): %v", err)
	}

	hashcash := hashcash.New(cfg.LeadingBytes, cfg.SaltSize, "")
	book := book.New(cfg.QuotesPath, log)

	handlerQuotes := quotes.New(book)
	handlerVerify := verify.New(hashcash)

	server := server.New(log)
	server.RegisterHandler("verify", handlerVerify)
	server.RegisterHandler("quote", handlerQuotes)

	server.Listen(cfg.Address)
}
