package zeus_server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type Server struct {
	E       *echo.Echo
	host    string
	port    string
	Name    string
	K8sUtil autok8s_core.K8Util
}

type Config struct {
	Host      string
	Port      string
	Name      string
	PGConnStr string
	K8sUtil   autok8s_core.K8Util
}

func NewZeusServer(cfg Config) Server {
	srv := Server{
		host:    cfg.Host,
		port:    cfg.Port,
		E:       InitBaseRoute(),
		K8sUtil: cfg.K8sUtil,
	}
	return srv
}

func InitBaseRoute() *echo.Echo {
	e := echo.New()
	e.Use(
		middleware.Recover(),
		middleware.Logger(),
	)
	return e
}

func (s *Server) Start() {
	address := fmt.Sprintf("%s:%s", s.host, s.port)

	// Start server
	go func() {
		if err := s.E.Start(address); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Start up failed, shutting down the server")
		}
	}()
	log.Info().Msgf("server listening on address %s", address)

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.E.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Start up failed, shutting down the server")
	}
}
