package exiter

import (
	"context"
	"github.com/dimitryshirokov/simple-app/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func HandleExit(ctx context.Context, s *server.Server) {
	cancelCtx, cancel := context.WithCancel(ctx)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()
	for {
		select {
		case <-cancelCtx.Done():
			s.Stop()
			return
		}
	}
}
