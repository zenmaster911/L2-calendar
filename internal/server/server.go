package server

import (
	"net/http"
	"time"

	"github.com/znmaster911/L2-calendar/internal/logger"
)

type Server struct {
	HttpsServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler, logger *logger.Logger) error {
	s.HttpsServer = &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	return s.HttpsServer.ListenAndServe()
}
