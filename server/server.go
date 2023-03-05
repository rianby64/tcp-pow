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
	processingUsers   *sync.WaitGroup
	processTimeout    time.Duration
	log               Logger
}

func (server *Server) handler(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("conn.Close(): %v", err)
		}
	}()

	done := make(chan struct{}, 1)

	go func() {
		defer func() {
			done <- struct{}{}
		}()

		for _, item := range server.registredHandlers {
			if err := item.handler.Handler(conn); err != nil {
				log.Printf("(%s) handler.Handler(conn): %v", item.name, err)

				return
			}
		}
	}()

	select {
	case <-done:
		log.Printf("done")

		return
	case <-time.After(server.processTimeout):
		log.Printf("timeout ocurred")

		return
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
	}
}
