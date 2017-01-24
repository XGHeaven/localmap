package server

import (
	"net"

	"github.com/xgheaven/localmap/logger"
)

var Option ServerOption

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
