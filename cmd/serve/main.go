package main

import (
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/catalinc/hashcash"
	"github.com/rianby64/tcp-pow/book"
	"github.com/rianby64/tcp-pow/handler/quotes"
	"github.com/rianby64/tcp-pow/handler/verify"
)

const (
	network = "tcp"
	address = ":9999"
)

type Logger interface {
	Printf(format string, v ...any)
}

type Handler interface {
	Handler(conn net.Conn) error
}

func handler(conn net.Conn, handlerQuotes, handlerVerify Handler, hashcash *hashcash.Hash, log Logger) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("conn.Close(): %v", err)
		}
	}()

	if err := handlerVerify.Handler(conn); err != nil {
		log.Printf("handlerVerify.Handler(conn): %v", err)

		return
	}

	if err := handlerQuotes.Handler(conn); err != nil {
		log.Printf("handlerQuotes.Handler(conn): %v", err)
	}
}

func main() {
	log := log.New(os.Stderr, "", log.LstdFlags)

	rand.Seed(time.Now().Unix())

	leadingBits := uint(20)
	saltLength := uint(8)
	extra := ""

	hashcash := hashcash.New(leadingBits, saltLength, extra)
	book := book.New("./quotes.json", log)

	handlerQuotes := quotes.New(book)
	handlerVerify := verify.New(hashcash)

	listener, err := net.Listen(network, address)

	if err != nil {
		log.Panicf("net.Listen(network=%s, address=%s): %v", network, address, err)
	}

	log.Printf("Listening at [%s]%s", network, address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("listener.Accept(): %v", err)
		}

		go handler(conn, handlerQuotes, handlerVerify, hashcash, log)
	}
}
