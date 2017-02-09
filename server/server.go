package server

import (
	"net"

	"github.com/xgheaven/localmap/logger"
)

var Option ServerOption

// Server server object
type Server struct {
	*ServerOption
	lsn     *net.TCPListener
	clients map[string]*ServerClient
}

func NewServer(option *ServerOption) *Server {
	server := &Server{
		ServerOption: option,
		clients:      make(map[string]*ServerClient),
	}
	return server
}

func Start() {
	listener, err := net.ListenTCP("tcp", &Option.Addr)

	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("listen at " + Option.Addr.String())

	for {
		c, err := listener.AcceptTCP()
		if err != nil {
			logger.Error(err)
			continue
		}
		logger.Debug("accept new client")
		go HandleConnect(c)
	}
}
