package server

import (
	"net"
	"sync"

	"github.com/pkg/errors"
)

func (server *Server) Listen(address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return errors.Wrapf(err, "net.Listen(network=%s, address=%s): %v", network, address, err)
	}

	server.processingUsers = &sync.WaitGroup{}
	server.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			return errors.Wrap(err, "listener.Accept()")
		}

		server.processingUsers.Add(1)

		go func() {
			defer server.processingUsers.Done()

			server.handler(conn)
		}()
	}
}
