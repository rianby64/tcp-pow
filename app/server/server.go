package server

import (
	"log"
	"net"
	"sync"
	"time"
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
	processTimeout    time.Duration
	log               Logger

	wgCurrentHandling *sync.WaitGroup
}

func (server *Server) handler(conn net.Conn) {
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

func New(log Logger, processTimeoutSecs int) *Server {
	return &Server{
		registredHandlers: make(registredHandlers, 0),
		processTimeout:    time.Duration(processTimeoutSecs) * time.Second,
		log:               log,

		wgCurrentHandling: &sync.WaitGroup{},
	}
}
