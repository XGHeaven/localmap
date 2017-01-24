package main

import (
	"errors"
	"flag"
	"io"
	"net"

	"github.com/xgheaven/localmap/client"
	"github.com/xgheaven/localmap/logger"
	"github.com/xgheaven/localmap/server"
)

func getMessage(conn net.Conn) {
	io.Copy(conn, conn)
}

var (
	isServer bool
	isClient bool
	sPort    int
	cPort    int
	sAddr    string
)

func init() {
	flag.BoolVar(&isServer, "server", false, "-server")
	flag.BoolVar(&isClient, "client", false, "-client")
	flag.IntVar(&sPort, "sport", 8000, "-sport=8000")
	flag.IntVar(&cPort, "cport", 8080, "-cport=8080")
	flag.StringVar(&sAddr, "addr", "127.0.0.1", "-addr=127.0.0.1")
	flag.Parse()
	if !isClient {
		isServer = true
		isClient = false
		err := checkServerOption()
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("start as server")
	} else {
		isServer = false
		isClient = true
		err := checkClientOption()
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("start as client")
	}
}

func checkServerOption() error {
	return nil
}

func checkClientOption() error {
	if sAddr == "" {
		return errors.New("empty server addr")
	}
	return nil
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			logger.Error(err)
		}
	}()
	if isServer {
		server.Option = server.ServerOption{
			Addr: net.TCPAddr{IP: net.IPv4zero, Port: sPort},
		}
		server.Start()
	}
	if isClient {
		client.Option = client.ClientOption{
			SAddr: net.TCPAddr{IP: net.IP{127, 0, 0, 1}, Port: sPort},
			CAddr: net.TCPAddr{IP: net.IP{127, 0, 0, 1}, Port: cPort},
		}
		client.Start()
	}
}
