package quotes

import (
	"encoding/json"
	"net"

	"github.com/pkg/errors"
	"github.com/rianby64/tcp-pow/app/models"
)

type Logger interface {
	Printf(format string, v ...any)
}

type Book interface {
	GetRandomQuote() *models.Quote
}

type Handler struct {
	book Book
}

func (handler *Handler) Handler(conn net.Conn) error {
	quote := handler.book.GetRandomQuote()
	data, err := json.Marshal(quote)

	if err != nil {
		return errors.Wrap(err, "json.Marshal(quote)")
	}

	if _, err := conn.Write(data); err != nil {
		return errors.Wrap(err, "conn.Write(data)")
	}

	return nil
}

func New(book Book) *Handler {
	return &Handler{
		book: book,
	}
}
