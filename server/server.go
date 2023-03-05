package server

import (
	"log"
	"net"
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

	log Logger
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

func (server *Server) Listen(address string) {
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

		go server.handler(conn)
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
