package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/wwq1988/leaf/pkg/conf"
	"github.com/wwq1988/leaf/pkg/service"
)

type Server struct {
	server  *http.Server
	service service.Service
	ln      net.Listener
}

func New(service service.Service) *Server {
	return &Server{service: service}
}

func (s *Server) Init() error {
	addr := conf.GetListenAddr()
	var err error
	s.ln, err = net.Listen("tcp", addr)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("addr:%s", addr))
	}
	s.server = &http.Server{Handler: s.service}
	return nil
}

func (s *Server) Start() error {
	if err := s.server.Serve(s.ln); err != http.ErrServerClosed {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Server) Stop() {
	s.server.Shutdown(context.Background())
	s.service.Stop()
}
