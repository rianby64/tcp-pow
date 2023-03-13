package server

import (
	"context"
	"net"
	"time"

	"github.com/pkg/errors"
)

func (server *Server) handleConn(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			server.log.Printf("conn.Close(): %v", err)
		}

		server.wgCurrentHandling.Done()
	}()

	deadline := time.Now().Add(server.processTimeout)

	if err := conn.SetDeadline(deadline); err != nil {
		server.log.Printf("conn.SetDeadline(deadline=%v)): %v", deadline, err)
	}

	server.handler(conn)
}

func (server *Server) infiniteLoop() error {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			return errors.Wrap(err, "server.listener.Accept()")
		}

		server.wgCurrentHandling.Add(1)

		go server.handleConn(conn)
	}
}

func (server *Server) Listen(ctx context.Context, address string) error {
	var lc net.ListenConfig

	listener, err := lc.Listen(ctx, network, address)
	if err != nil {
		return errors.Wrapf(err, "net.Listen(network=%s, address=%s)", network, address)
	}

	server.listener = listener

	if err := server.infiniteLoop(); err != nil {
		return errors.Wrap(err, "err := server.infiniteLoop()")
	}

	return nil
}
