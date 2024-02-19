package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/dimitryshirokov/simple-app/internal/config"
	"github.com/dimitryshirokov/simple-app/internal/logger"
	"github.com/dimitryshirokov/simple-app/internal/server/handler"
	"net/http"
	"sync"
	"time"
)

func NewServer(ctx context.Context, c *config.Config, handlers map[string]handler.Handler) *Server {
	return &Server{
		ctx:      ctx,
		c:        c,
		handlers: handlers,
	}
}

type Server struct {
	ctx      context.Context
	c        *config.Config
	handlers map[string]handler.Handler
	srv      *http.Server
}

func (s *Server) Run() *sync.WaitGroup {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.c.HttpPort),
		ReadHeaderTimeout: 10 * time.Second,
	}
	for route, h := range s.handlers {
		http.HandleFunc(route, h.Handle)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.LogError("can't create http server", nil, err)
		}
	}()
	s.srv = srv
	return wg
}

func (s *Server) Stop() {
	timeoutCtx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	err := s.srv.Shutdown(timeoutCtx)
	if err != nil {
		logger.LogError("can't stop http server", nil, err)
		panic(err)
	}
}
