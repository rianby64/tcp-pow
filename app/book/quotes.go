package book

import (
	"encoding/json"
	"math/rand"
	"os"

	"github.com/rianby64/tcp-pow/app/models"
)

type Logger interface {
	Panicf(format string, v ...any)
}

type Book struct {
	quotes models.Quotes
	size   int
}

func (book *Book) GetRandomQuote() *models.Quote {
	randomIndex := rand.Intn(book.size)

	return book.quotes[randomIndex]
}

func New(path string, log Logger) *Book {
	var quotes models.Quotes

	if data, err := os.ReadFile(path); err != nil {
		log.Panicf("os.ReadFile(path=%s): %v", path, err)
	} else if err := json.Unmarshal(data, &quotes); err != nil {
		log.Panicf("json.Unmarshal(data, &quotes) (path=%s): %v", path, err)
	}

	size := len(quotes)
	if size == 0 {
		log.Panicf("size == 0, %v", models.ErrEmpty)
	}

	return &Book{
		quotes: quotes,
		size:   size,
	}
}
