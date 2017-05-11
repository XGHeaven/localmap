package server

import (
	"errors"
	"net"

	"github.com/xgheaven/localmap/connect"
	"github.com/xgheaven/localmap/logger"
)

// Start start server to listen port
func (server *Server) Start() error {
	if server.lsn != nil {
		return errors.New("server has been started")
	}

	listener, err := net.ListenTCP("tcp", &server.Addr)

	if err != nil {
		logger.Error(err)
		return err
	}

	server.lsn = listener
	logger.Info("server started at", server.Addr.String())

	server.AcceptClient()

	return nil
}

func (server *Server) AcceptClient() {
	listener := server.lsn
	for {
		clientConn, err := listener.AcceptTCP()
		logger.Info("accept a new client from", clientConn.RemoteAddr().String())

		if err != nil {
			logger.Warning("accept new client error")
			continue
		}

		cli := NewServerClient(server, (*connect.TCPConnect)(clientConn))
		//server.clients[cli.UUID.String()] = cli

		go cli.Start()
	}
}
