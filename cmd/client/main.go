package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/catalinc/hashcash"
	"github.com/pkg/errors"
)

const (
	sizeReadBytes = 31
)

type Config struct {
	Address     string `env:"ADDRESS,required"`
	LeadingBits uint   `env:"LEADING_BITS,required"`
	SaltSize    uint   `env:"SALT_SIZE,required"`
	KeySize     int    `env:"KEY_SIZE" envDefault:"64"`
}

func createLogAndConfig() (*log.Logger, *Config) {
	cfg := &Config{}
	log := log.New(os.Stderr, "client-tcp-pow", log.LstdFlags)

	if err := env.Parse(cfg); err != nil {
		log.Panicf("env.Parse(&cfg): %v", err)
	}

	log.Printf("ADDRESS=%v", cfg.Address)
	log.Printf("LEADING_BITS=%v", cfg.LeadingBits)
	log.Printf("SALT_SIZE=%v", cfg.SaltSize)

	return log, cfg
}

func Generate(leadingBits, saltLength uint, data string) (string, error) {
	extra := ""
	h := hashcash.New(leadingBits, saltLength, extra)

	// Mint a new stamp
	stamp, err := h.Mint(data)
	if err != nil {
		return "", errors.Wrap(err, "stamp, err := h.Mint(string(data))")
	}

	return stamp, nil
}

func readMsg(conn net.Conn) ([]byte, error) {
	var msg []byte

	chunk := make([]byte, sizeReadBytes)

	for {
		n, err := conn.Read(chunk)
		if errors.Is(err, io.EOF) {
			msg = append(msg, chunk[:n]...)

			break
		}

		if err != nil {
			return nil, errors.Wrap(err, "n, err := conn.Read(chunk)")
		}

		msg = append(msg, chunk[:n]...)
	}

	return msg, nil
}

func main() {
	log, cfg := createLogAndConfig()
	conn, err := net.Dial("tcp", cfg.Address)

	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	puzzle := make([]byte, cfg.KeySize)
	if n, err := conn.Read(puzzle); err != nil {
		log.Panic(err)
	} else {
		puzzle = puzzle[:n]
	}

	if err != nil {
		log.Panic(err)
	}

	puzzleResolved, err := Generate(cfg.LeadingBits, cfg.SaltSize, string(puzzle))
	if err != nil {
		log.Panic(err)
	}

	if _, err := conn.Write([]byte(puzzleResolved)); err != nil {
		log.Panic(err)
	}

	quote, err := readMsg(conn)
	if err != nil {
		log.Panic(err)
	}

	log.Println(string(quote))
}
