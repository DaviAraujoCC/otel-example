package server

import (
	"context"
	"net/http"

	"products-api/server/router"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
)

type Server interface {
	Listen()
	Shutdown()
}

type server struct {
	router       *mux.Router
	httpServer   *http.Server
	shutdownChan chan struct{}
}

func NewServer() Server {

	s := &server{
		router:       router.New(),
		shutdownChan: make(chan struct{}, 1),
	}

	setupServer(s)

	return s
}

func setupServer(s *server) {
	neg := negroni.Classic()
	neg.UseHandler(s.router)

	s.httpServer = &http.Server{
		Addr:    ":" + viper.GetString("HTTP_PORT"),
		Handler: neg,
	}
}

func (s server) Listen() {

	ctx := context.Background()

	go func() {
		<-s.shutdownChan
		logrus.Warn("shutting down server")

		logrus.Warn("exiting...")
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			logrus.Fatalf("Error while shutting down server: %v", err.Error())
		}
	}()

	go func() {

		if err := s.httpServer.ListenAndServe(); err != nil {
			logrus.WithError(err).Warn("server shutdown")
			s.Shutdown()
		}

	}()

}

func (s *server) Shutdown() {
	s.shutdownChan <- struct{}{}
}