package main

import (
	"log"

	"github.com/rianby64/tcp-pow/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Panic(err)
	}
}
