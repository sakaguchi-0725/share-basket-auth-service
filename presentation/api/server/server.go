package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

type Server struct {
	addr string
	*echo.Echo
}

func New(addr string) *Server {
	e := echo.New()
	return &Server{addr, e}
}

func (s *Server) Run() error {
	server := &http.Server{Addr: s.addr}

	go func() {
		if err := s.StartServer(server); err != nil {
			log.Fatalf("Could not start the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.Shutdown(ctx)
}
