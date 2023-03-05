package main

import (
	"io"
	"log"
	"net"

	"github.com/catalinc/hashcash"
	"github.com/pkg/errors"
)

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

func main() {
	client, err := net.Dial("tcp", ":9999")

	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Println(err)
		}
	}()

	msg := make([]byte, 128)
	if n, err := client.Read(msg); err != nil {
		log.Panic(err)
	} else {
		msg = msg[:n]
	}

	puzzleResolved, err := Generate(20, 8, string(msg))
	if err != nil {
		log.Panic(err)
	}

	if _, err := client.Write([]byte(puzzleResolved)); err != nil {
		log.Panic(err)
	}

	quote := []byte{}

	for {
		if n, err := client.Read(msg); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Panic(err)
		} else {
			msg = msg[:n]
		}

		quote = append(quote, msg...)
	}

	log.Println(string(quote))
}
