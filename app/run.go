package app

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/catalinc/hashcash"
	"github.com/pkg/errors"
	"github.com/rianby64/tcp-pow/app/book"
	"github.com/rianby64/tcp-pow/app/config"
	"github.com/rianby64/tcp-pow/app/handler/quotes"
	"github.com/rianby64/tcp-pow/app/handler/verify"
	"github.com/rianby64/tcp-pow/app/server"
)

type Shutdown interface {
	Shutdown(ctx context.Context) error
}

func waitSignalShutdown() <-chan os.Signal {
	signalShutdown := make(chan os.Signal, 1)

	signal.Notify(signalShutdown, syscall.SIGINT, syscall.SIGTERM)

	return signalShutdown
}

func shutdown(log *log.Logger, shutdownTimeout time.Duration, items ...Shutdown) {
	log.Printf("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	defer cancel()

	for _, item := range items {
		if err := item.Shutdown(ctx); err != nil {
			log.Printf("err := item.Shutdown(ctx): %v", err)
		}
	}

	log.Printf("shutdown completed")
}

func listen(server *server.Server, log *log.Logger, address string, shutdownTimeout time.Duration) <-chan error {
	log.Printf("Listening at %s", address)

	ctx := context.Background()
	errChan := make(chan error, 1)

	go func() {
		err := server.Listen(ctx, address)

		if errors.Is(err, net.ErrClosed) {
			log.Printf("Listener closed")

			errChan <- nil
		}

		if err != nil {
			errChan <- errors.Wrapf(err, "server.Listen(%s)", address)
		}
	}()

	return errChan
}

func Run() {
	log, cfg := config.CreateLogAndConfig()

	hashcash := hashcash.New(cfg.LeadingBits, cfg.SaltSize, "")
	book := book.New(cfg.QuotesPath, log)

	handlerQuotes := quotes.New(book)
	handlerVerify := verify.New(hashcash, cfg.KeySize, cfg.LeadingBits, cfg.SaltSize)

	server := server.New(log, cfg.ProcessTimeoutSecs)
	server.RegisterHandler("verify", handlerVerify)
	server.RegisterHandler("quote", handlerQuotes)

	shutdownTimeout := time.Duration(cfg.ShutdownTimeoutSecs) * time.Second

	defer shutdown(log, shutdownTimeout, server)

	select {
	case err := <-listen(server, log, cfg.Address, shutdownTimeout):
		if err != nil {
			log.Panicf("listener: %v", err)
		}
	case <-waitSignalShutdown():
		log.Printf("waitSignalShutdown received")
	}
}
