package server

import (
	"log"
	"net"
	"sync"
)

const (
	network = "tcp"
)

type Handler interface {
	Handler(conn net.Conn) error
}

type registredHandler struct {
	handler Handler
	name    string
}

type registredHandlers []*registredHandler

type Logger interface {
	Printf(format string, v ...any)
	Panicf(format string, v ...any)
}

type Server struct {
	registredHandlers registredHandlers
	listener          net.Listener
	processingUsers   *sync.WaitGroup
	log               Logger
}

func (server *Server) handler(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("conn.Close(): %v", err)
		}
	}()

	for _, item := range server.registredHandlers {
		if err := item.handler.Handler(conn); err != nil {
			log.Printf("(%s) handler.Handler(conn): %v", item.name, err)

			return
		}
	}
}

func (server *Server) RegisterHandler(name string, handler Handler) {
	server.registredHandlers = append(server.registredHandlers, &registredHandler{
		name:    name,
		handler: handler,
	})
}

func New(log Logger) *Server {
	return &Server{
		registredHandlers: make(registredHandlers, 0),
		log:               log,
	}
}
