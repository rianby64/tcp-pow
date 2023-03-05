package main

import (
	"github.com/catalinc/hashcash"
	"github.com/rianby64/tcp-pow/book"
	"github.com/rianby64/tcp-pow/handler/quotes"
	"github.com/rianby64/tcp-pow/handler/verify"
	"github.com/rianby64/tcp-pow/server"
)

func main() {
	log, cfg := createLogAndConfig()
	hashcash := hashcash.New(cfg.LeadingBits, cfg.SaltSize, "")
	book := book.New(cfg.QuotesPath, log)

	handlerQuotes := quotes.New(book)
	handlerVerify := verify.New(hashcash, cfg.KeySize, cfg.LeadingBits, cfg.SaltSize)

	server := server.New(log)
	server.RegisterHandler("verify", handlerVerify)
	server.RegisterHandler("quote", handlerQuotes)

	go listen(server, log, cfg.Address)

	waitSignalShutdown()

	shutdown(log, server)
}
