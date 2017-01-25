package main

import (
	"errors"
	"flag"
	"io"
	"net"
	"os"

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
	sAddr    net.IP
	showHelp bool
	_sAddr   string
)

func init() {
	flag.BoolVar(&isServer, "server", false, "start as server mode")
	flag.BoolVar(&isClient, "client", false, "start as client mode")
	flag.IntVar(&sPort, "sport", 8000, "set server port to connect (only client mode)")
	flag.IntVar(&cPort, "cport", -1, "set client port to connect (only client mode)")
	flag.StringVar(&_sAddr, "addr", "127.0.0.1", "where server address to connect, support ip, domain")
	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

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
	if addr, err := net.ResolveIPAddr("ip", _sAddr); err != nil {
		return errors.New("wrong server address, please use right address")
	} else {
		sAddr = addr.IP
	}
	if cPort == -1 {
		return errors.New("please specify which port to connect for client")
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
			SAddr: net.TCPAddr{IP: sAddr, Port: sPort},
			CAddr: net.TCPAddr{IP: net.IP{127, 0, 0, 1}, Port: cPort},
		}
		client.Start()
	}
}
