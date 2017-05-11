package server

import (
	"net"
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
