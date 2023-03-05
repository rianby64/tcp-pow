package server

import (
	"context"
)

func (server *Server) Shutdown(ctx context.Context) (err error) {
	if server.listener == nil {
		return nil
	}

	err = server.listener.Close()
	allDone := make(chan struct{}, 1)

	go func() {
		server.processingUsers.Wait()

		allDone <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return
	case <-allDone:
		return
	}
}
