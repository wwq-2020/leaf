package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wwq1988/leaf/pkg/generator"
	"go.uber.org/zap"
)

type Service interface {
	http.Handler
	Stop()
}

type service struct {
	logger    *zap.Logger
	generator generator.Generator
	http.Handler
}

func New(generator generator.Generator, logger *zap.Logger) Service {
	service := &service{logger: logger, generator: generator}

	mux := mux.NewRouter().StrictSlash(true).SkipClean(true)
	sub := mux.PathPrefix("/api").Subrouter().StrictSlash(true).SkipClean(true)
	sub.HandleFunc("/get/{key}", service.get)
	service.Handler = mux
	return service
}

func (s *service) Stop() {
	s.generator.Stop()
}
