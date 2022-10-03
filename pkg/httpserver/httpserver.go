package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/m1g88/go-api-scaffold/pkg/config"
)

type Config struct {
	SrvPort string `envconfig:"SERVER_PORT" default:"8080"`
	ENV     string `envconfig:"ENV" default:"local"`
}

type Option func(*server)

func WithPort(port int) Option {
	return func(s *server) {
		s.port = port
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.timeout = timeout
	}
}

type server struct {
	port    int
	timeout time.Duration

	Handler http.Handler

	server *http.Server
}

func NewServer(opts ...Option) *server {
	var cfg Config
	config.MustProcess("", &cfg)
	s := new(server)

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func NewGinServer(opts ...Option) *server {
	var cfg Config
	config.MustProcess("", &cfg)

	s := new(server)
	s.port, _ = strconv.Atoi(cfg.SrvPort)

	if cfg.ENV == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (srv *server) Run() {
	srv.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", srv.port),
		Handler: srv.Handler,
	}

	go func(server *http.Server) {
		log.Printf("Listen on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error while listen and serve: %v", err)
		}
	}(srv.server)

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt, syscall.SIGTERM)
	<-wait

	log.Printf("Shutting down http server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := srv.server.Shutdown(ctx); err != nil {
		log.Printf("Cannot shutdown server: %v", err)
	}

	cancel()
}
