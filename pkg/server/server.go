package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Server defines a struct for server
type Server struct {
	log           *logrus.Logger
	Router        *gin.Engine
	httpServer    *http.Server
	configuration Configuration
	ErrCh         chan error
}

// NewServer initializes a server
func NewServer(log *logrus.Logger, configuration Configuration) *Server {
	router := gin.Default()

	s := &Server{
		log:           log,
		configuration: configuration,
	}

	router.Use(
		s.tracing(),
	)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", configuration.Host, configuration.Port),
		Handler: router,
	}

	log.Info(s.httpServer.Addr)

	s.Router = router

	return s
}

// Run when called starts the server
func (s *Server) Run(ctx context.Context) <-chan error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			errM := fmt.Sprintf("unexpected error while running server %v", err.Error())
			s.ErrCh <- errors.New(errM)
		}

		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		s.log.Infof("Shutdown Server ...")

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.log.Fatalf("Server forced to shutdown: %v", err)
		}

		s.log.Infof("Server exiting")
	}()

	return s.ErrCh
}
